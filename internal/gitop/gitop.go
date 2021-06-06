package gitop

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
)

// TODO move to configuration
// FIXME we need proper path
var gitConfig = map[string]map[string]string{
	"mw": {
		"url":        "ssh://mnews@mnews.tw@source.developers.google.com:2022/p/mirrormedia-1470651750304/r/helm",
		"branch":     "master",
		"sshKeyPath": "/Users/chiu/dev/mtv/major-tom-go/configs/identity",
		"sshKeyUser": "mnews@mnews.tw",
	},
	"tv": {
		"url":        "ssh://source.developers.google.com:2022/p/mirror-tv-275709/r/helm",
		"branch":     "master",
		"sshKeyPath": "/Users/chiu/dev/mtv/major-tom-go/configs/identity",
		"sshKeyUser": "mnews@mnews.tw",
	},
	"readr": {
		"url":        "ssh://mnews@mnews.tw@source.developers.google.com:2022/p/mirrormedia-1470651750304/r/helm",
		"branch":     "master",
		"sshKeyPath": "/Users/chiu/dev/mtv/major-tom-go/configs/identity",
		"sshKeyUser": "mnews@mnews.tw",
	},
}

type Repository struct {
	r      *git.Repository
	once   *sync.Once
	config map[string]string
}

var mw, tv, readr = &Repository{
	r:      nil,
	once:   &sync.Once{},
	config: gitConfig["mw"],
}, &Repository{
	r:      nil,
	once:   &sync.Once{},
	config: gitConfig["tv"],
}, &Repository{
	r:      nil,
	once:   &sync.Once{},
	config: gitConfig["readr"],
}

// GetFile will return a File interface with read and write permission
func (repo *Repository) GetFile(filenamePath string) (billy.File, error) {
	worktree, err := repo.r.Worktree()
	if err != nil {
		return nil, nil
	}
	f, err := worktree.Filesystem.OpenFile(filenamePath, 0666, os.ModePerm)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))

	return f, err
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
			Name:  name,
			Email: email,
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

func GetRepository(project string) (r *Repository, err error) {

	// Get the singleton repository according to the project
	return getRepository(project)
}

func getRepository(project string) (repo *Repository, err error) {
	switch project {
	case "mw":
		repo = mw
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
			URL:           config["url"],
			ReferenceName: plumbing.NewBranchReferenceName(config["branch"]),
			SingleBranch:  true,
			Auth:          sshMethod,
			// Set depth to 1 because we only need the HEAD
			Depth: 1,
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
	return repo, err
}