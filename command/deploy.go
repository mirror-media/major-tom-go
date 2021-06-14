package command

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	gootkitconfig "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/mirror-media/major-tom-go/config"
	"github.com/mirror-media/major-tom-go/gitop"
	mjcontext "github.com/mirror-media/major-tom-go/internal/context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	once       sync.Once
	gitConfigs map[config.Repository]config.GitConfig
}

var DeployWorker deployWorker

func (w *deployWorker) Init(gitConfigs map[config.Repository]config.GitConfig) {
	w.gitConfigs = gitConfigs
	go w.once.Do(func() {
		for {
			deployment := <-deployChannel
			deploy(deployment.ctx, deployment.clusterConfigs, w.gitConfigs, deployment.textParts, deployment.caller)
		}
	})
}

type operation struct {
	keyPath string
	value   interface{}
}

// deploy deploy certain configuration to a service. textParts in interpreted as [project, stage, service, ...cfg:value]
func deploy(ctx context.Context, clusterConfigs config.K8S, gitConfigs map[config.Repository]config.GitConfig, textParts []string, caller string) {
	var messages = []string{}
	var err error
	ch := ctx.Value(mjcontext.ResponseChannel).(chan response)
	switch len(textParts) {
	// Deploy needs cfg:arg to operation on
	case 0, 1, 2, 3:
		messages = append([]string{"call help"}, messages...)
		err = errors.New("Please provide more details...")
	// textParts contains the config to deploy
	default:
		project := textParts[0]
		repo, errRepo := gitop.GetRepository(textParts[0], gitConfigs)
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
		operations := map[string]operation{}
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
				operations[cfg] = operation{
					keyPath: config.path,
					value:   value,
				}
			case "int":
				var v int64
				v, err = strconv.ParseInt(value, 10, 0)
				operations[cfg] = operation{
					keyPath: config.path,
					value:   v,
				}
			case "bool":
				var v bool
				v, err = strconv.ParseBool(value)
				operations[cfg] = operation{
					keyPath: config.path,
					value:   v,
				}
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
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("getting %s has error", valuesFilePath))
			logrus.Warn(err)
			ch <- response{
				Messages: messages,
				Error:    err,
			}
			return
		}

		// operation starts here. worktree needs to be cleaned if disaster happens
		b, err := io.ReadAll(f)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("reading %s has error", valuesFilePath))
			logrus.Warn(err)
			ch <- response{
				Messages: messages,
				Error:    err,
			}
			_ = hardResetFn()
			return
		}
		valueConfig := gootkitconfig.New(valuesFilePath)
		valueConfig.AddDriver(yaml.Driver)
		err = valueConfig.LoadStrings(gootkitconfig.Yaml, string(b))
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("loading YAML from %s has error", valuesFilePath))
			logrus.Warn(err)
			ch <- response{
				Messages: messages,
				Error:    err,
			}
			_ = hardResetFn()
			return
		}

		for _, operation := range operations {
			valueConfig.Set(operation.keyPath, operation.value, true)
		}

		err = f.Truncate(0)
		f.Close()
		if err != nil {
			ch <- response{
				Messages: messages,
				Error:    errors.Wrap(err, fmt.Sprintf("clean file before writing to %s has error", valuesFilePath)),
			}
			_ = hardResetFn()
			return
		}

		f, _ = repo.GetFile(valuesFilePath)
		valueConfig.DumpTo(f, gootkitconfig.Yaml)
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
