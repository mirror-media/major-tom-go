package command

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mirror-media/major-tom-go/v2/config"
	mjcontext "github.com/mirror-media/major-tom-go/v2/internal/context"
	"github.com/mirror-media/major-tom-go/v2/internal/gitop"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type yamlCFG struct {
	path      string
	valueType string
}

var valuesConfig = map[string]yamlCFG{
	"image":       yamlCFG{path: "image.tag", valueType: "string"},
	"autoScaling": yamlCFG{path: "autoscaling.enabled", valueType: "bool"},
	"pods":        yamlCFG{path: "replicacount", valueType: "int"},
	"maxPods":     yamlCFG{path: "autoscaling.maxReplicas", valueType: "int"},
	"minPods":     yamlCFG{path: "autoscaling.minReplicas", valueType: "int"},
}

type Deployment struct {
	ctx            context.Context
	clusterConfigs config.K8S
	textParts      []string
	caller         string
}

var deployChannel = make(chan Deployment, 64)

// Deploy certain configuration to a service. textParts in interpreted as [project, stage, service, ...cfg:arg]
func Deploy(ctx context.Context, clusterConfigs config.K8S, textParts []string, caller string) (messages []string, err error) {
	timeout := 5 * time.Minute
	ch := make(chan response)
	newCtx := context.WithValue(ctx, mjcontext.ResponseChannel, ch)
	newCtx, cancelFn := context.WithTimeout(newCtx, timeout)
	defer cancelFn()
	deployChannel <- Deployment{
		ctx:            newCtx,
		clusterConfigs: clusterConfigs,
		textParts:      textParts,
		caller:         caller,
	}

	select {
	case commandResponse := <-ch:
		return commandResponse.Messages, commandResponse.Error
	case <-newCtx.Done():
		return nil, errors.Errorf("\"%s\" command has timeouted(%f)", strings.Join(textParts, " "), timeout.Minutes())
	}
}

type deployWorker struct {
	once sync.Once
}

var DeployWorker deployWorker

func (w *deployWorker) Init() {
	go w.once.Do(func() {
		for {
			deployment := <-deployChannel
			deploy(deployment.ctx, deployment.clusterConfigs, deployment.textParts, deployment.caller)
		}
	})
}

// deploy deploy certain configuration to a service. textParts in interpreted as [project, stage, service, ...cfg:arg]
func deploy(ctx context.Context, clusterConfigs config.K8S, textParts []string, caller string) {
	var messages []string
	var err error
	ch := ctx.Value(mjcontext.ResponseChannel).(chan response)
	switch len(textParts) {
	// Deploy needs cfg:arg to operation on
	case 0, 1, 2, 3:
		messages = []string{"call help"}
		err = errors.New("Please provide more details...")
	// textParts contains the config to deploy
	default:
		project := textParts[0]
		repo, errRepo := gitop.GetRepository(textParts[0])
		if errRepo != nil {
			err = errors.Wrap(errRepo, fmt.Sprintf("getting repository for project(%s) has error", project))
			break
		}
		hash, errHash := repo.GetHeadHash()
		if errHash != nil {
			err = errors.Wrap(errHash, fmt.Sprintf("getting head hash of repo for project(%s) has error", project))
			break
		}
		hardResetFn := hardReset(repo, hash)
		stage := textParts[1]
		service := textParts[2]
		args := textParts[3:]
		operations := map[string]interface{}{}
		// Validate the arguments
		for _, arg := range args {
			split := strings.Split(arg, ":")
			if len(split) != 2 {
				err = errors.Errorf("\"%s\" is malformatted", arg)
				break
			}
			cfg, value := split[0], split[1]
			config, ok := valuesConfig[cfg]
			if !ok {
				err = errors.Errorf("config(%s) is no supported", cfg)
				break
			}
			var t string
			switch t = config.valueType; t {
			case "string":
				operations[config.path] = value
			case "int":
				operations[config.path], err = strconv.ParseInt(value, 10, 0)
			case "bool":
				operations[config.path], err = strconv.ParseBool(value)
			}
			if err != nil {
				err = errors.Wrap(err, fmt.Sprintf("value type of %s is %s, but %s is provided", cfg, config.valueType, value))
				break
			}
		}
		// validation result
		if err != nil {
			ch <- response{
				Messages: messages,
				Error:    err,
			}
			return
		}

		valuesFilePath := fmt.Sprintf("%s/values-%s.yaml", service, stage)
		f, err := repo.GetFile(valuesFilePath)
		defer f.Close()
		if err != nil {
			ch <- response{
				Messages: messages,
				Error:    errors.Wrap(err, fmt.Sprintf("getting %s has error", valuesFilePath)),
			}
			return
		}

		// operation starts here. worktree needs to be cleaned if disaster happens
		viper := viper.New()
		viper.SetConfigType("yaml")
		err = viper.ReadConfig(f)
		if err != nil {
			ch <- response{
				Messages: messages,
				Error:    errors.Wrap(err, fmt.Sprintf("parsing %s has error", valuesFilePath)),
			}
			_ = hardResetFn()
			return
		}
		for key, value := range operations {
			viper.Set(key, value)
		}
		c := viper.AllSettings()

		b, err := yaml.Marshal(c)
		if err != nil {
			ch <- response{
				Messages: messages,
				Error:    errors.Wrap(err, fmt.Sprintf("marshalling %s has error", valuesFilePath)),
			}
			_ = hardResetFn()
			return
		}
		// f.Close() is deferred already
		err = f.Truncate(0)
		if err != nil {
			ch <- response{
				Messages: messages,
				Error:    errors.Wrap(err, fmt.Sprintf("clean file before writing to %s has error", valuesFilePath)),
			}
			_ = hardResetFn()
			return
		}
		err = func() error {
			f, _ = repo.GetFile(valuesFilePath)
			defer f.Close()
			_, err = f.Write(b)
			return err
		}()
		if err != nil {
			ch <- response{
				Messages: messages,
				Error:    errors.Wrap(err, fmt.Sprintf("writing to %s has error", valuesFilePath)),
			}
			_ = hardResetFn()
			return
		}

		// command operation finished
		// now git operations starts

		repo.AddFile(valuesFilePath)
		if err != nil {
			ch <- response{
				Messages: messages,
				Error:    errors.Wrap(err, fmt.Sprintf("adding %s to staging area has error", valuesFilePath)),
			}
			_ = hardResetFn()
			return
		}

		var body string
		sortedOperations := make([]string, 0, len(operations))
		for path := range operations {
			sortedOperations = append(sortedOperations, path)
		}
		sort.Strings(sortedOperations)
		for _, path := range sortedOperations {
			body += fmt.Sprintf("Set %s to %v\n", path, operations[path])
		}
		message := fmt.Sprintf("release(%s/%s): released by %s\n\n%s", service, stage, caller, body)

		err = repo.Commit(valuesFilePath, caller, message)
		if err != nil {
			ch <- response{
				Messages: append([]string{"this operation failed"}, messages...),
				Error:    errors.Wrap(err, fmt.Sprintf("commits for project(%s) has error", project)),
			}
			_ = hardResetFn()
			return
		}
		err = repo.Push()
		if err != nil {
			ch <- response{
				Messages: append([]string{"this operation failed"}, messages...),
				Error:    errors.Wrap(err, fmt.Sprintf("push commits for project(%s) has error", project)),
			}
			_ = hardResetFn()
			return
		}
		messages = strings.Split(message, "\n")
	}
	ch <- response{
		Messages: messages,
		Error:    err,
	}
}

func hardReset(repository *gitop.Repository, commit plumbing.Hash) (hardResetFN func() error) {
	return func() error { return repository.HardResetToCommit(commit) }
}
