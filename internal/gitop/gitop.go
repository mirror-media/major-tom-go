package gitop

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TODO move to configuration
// FIXME we need proper path
var gitConfig = map[string]map[string]string{
	"mm": {
		"branch":     "master",
		"sshKeyPath": "/Users/chiu/dev/mtv/major-tom-go/configs/identity",
		"sshKeyUser": "mnews@mnews.tw",
		"url":        "ssh://mnews@mnews.tw@source.developers.google.com:2022/p/mirrormedia-1470651750304/r/helm",
	},
	"tv": {
		"branch":     "master",
		"sshKeyPath": "/Users/chiu/dev/mtv/major-tom-go/configs/identity",
		"sshKeyUser": "mnews@mnews.tw",
		"url":        "ssh://source.developers.google.com:2022/p/mirror-tv-275709/r/helm",
	},
	"readr": {
		"branch":     "master",
		"sshKeyPath": "/Users/chiu/dev/mtv/major-tom-go/configs/identity",
		"sshKeyUser": "mnews@mnews.tw",
		"url":        "ssh://mnews@mnews.tw@source.developers.google.com:2022/p/mirrormedia-1470651750304/r/helm",
	},
}

type Repository struct {
	config map[string]string
	once   *sync.Once
	r      *git.Repository
}

var mm, tv, readr = &Repository{
	config: gitConfig["mm"],
	once:   &sync.Once{},
	r:      nil,
}, &Repository{
	config: gitConfig["tv"],
	once:   &sync.Once{},
	r:      nil,
}, &Repository{
	config: gitConfig["readr"],
	once:   &sync.Once{},
	r:      nil,
}

// GetFile will return an io.ReadWriter with read and write permission
func (repo *Repository) GetFile(filenamePath string) (io.ReadWriter, error) {
	worktree, err := repo.r.Worktree()
	if err != nil {
		return nil, nil
	}
	f, err := worktree.Filesystem.OpenFile(filenamePath, 0666, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return f, err
}

// AddFiles add the file to the staging area of worktree
func (repo *Repository) AddFiles(filenamePath string) error {
	worktree, err := repo.r.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(filenamePath)

	logrus.Infof("$s is added to the staging area", filenamePath)

	return err
}

// Commit with username as slack caller name annotated by (Major Tom)
func (repo *Repository) Commit(filename, caller, message string) error {
	// TODO extract email and bot name as configuration
	return commit(repo.r, filename, fmt.Sprintf("%s(%s)", caller, "Major Tom"), "mnews@mnews.tw", message)
}

func commit(r *git.Repository, filename, name, email, message string) error {

	worktree, err := r.Worktree()
	if err != nil {
		return nil
	}

	commit, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Email: email,
			Name:  name,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	obj, err := r.CommitObject(commit)
	if err != nil {
		return err
	}

	fmt.Println(obj)
	return nil
}

func (repo *Repository) Pull() error {
	worktree, err := repo.r.Worktree()
	if err != nil {
		return err
	}
	return worktree.Pull(&git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(repo.config["branch"]),
		RemoteName:    "origin",
		SingleBranch:  true,
	})
}

func (repo *Repository) Push() error {
	return repo.r.Push(&git.PushOptions{})
}

func GetRepository(project string) (r *Repository, err error) {

	// Get the singleton repository according to the project
	return getRepository(project)
}

func getRepository(project string) (repo *Repository, err error) {
	switch project {
	case "mm":
		repo = mm
	case "tv":
		repo = tv
	case "readr":
		repo = readr
	default:
		return nil, errors.New("wrong project")
	}

	// Init git repo
	repo.once.Do(func() {
		// Get the config according to the project
		config := repo.config
		key, errRead := os.ReadFile(config["sshKeyPath"])
		if errRead != nil {
			err = errRead
			err = errors.Wrap(errRead, "reading ssh key failed")
			return
		}
		sshMethod, errSSH := ssh.NewPublicKeys(config["sshKeyUser"], key, "")
		if errSSH != nil {
			err = errors.Wrap(errSSH, "creating sshMethod from key failed")
			return
		}
		opt := git.CloneOptions{
			Auth: sshMethod,
			// Set depth to 1 because we only need the HEAD
			Depth:         1,
			ReferenceName: plumbing.NewBranchReferenceName(config["branch"]),
			SingleBranch:  true,
			URL:           config["url"],
		}
		newGitRepo, errGitRepo := cloneGitRepo(opt)
		if errGitRepo != nil {
			// Reset Once
			repo.once = &sync.Once{}
			err = errors.Wrap(errGitRepo, "cloning git repo failed")
		} else {
			repo.r = newGitRepo
		}
	})
	if err == nil {
		err = repo.Pull()
	}
	return repo, err
}
