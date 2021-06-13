package gitop

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mirror-media/major-tom-go/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	authMethod ssh.AuthMethod
	config     *config.GitConfig
	name       string
	once       *sync.Once
	r          *git.Repository
	locker     *sync.Mutex
}

var tv = &Repository{
	config: nil,
	name:   "mirror tv helm repo",
	once:   &sync.Once{},
	r:      nil,
	locker: &sync.Mutex{},
}

// var mm, tv, readr = &Repository{
// 	config: nil,
// 	name:   "mirror weekly helm repo",
// 	once:   &sync.Once{},
// 	r:      nil,
// 	locker: &sync.Mutex{},
// }, &Repository{
// 	config: nil,
// 	name:   "mirror tv helm repo",
// 	once:   &sync.Once{},
// 	r:      nil,
// 	locker: &sync.Mutex{},
// }, &Repository{
// 	config: nil,
// 	name:   "readr helm repo",
// 	once:   &sync.Once{},
// 	r:      nil,
// 	locker: &sync.Mutex{},
// }

// GetFile will return an billy.Filewith read and write permission
func (repo *Repository) GetFile(filenamePath string) (billy.File, error) {
	repo.locker.Lock()
	defer repo.locker.Unlock()
	worktree, err := repo.r.Worktree()
	if err != nil {
		return nil, err
	}
	f, err := worktree.Filesystem.OpenFile(filenamePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return f, err
}

// AddFile add the file to the staging area of worktree
func (repo *Repository) AddFile(filenamePath string) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()
	worktree, err := repo.r.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(filenamePath)

	logrus.Infof("%s is added to the staging area", filenamePath)

	return err
}

// GetHeadHash hard reset the worktree to the commit to clear changes
func (repo *Repository) GetHeadHash() (plumbing.Hash, error) {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	head, err := repo.r.Head()
	if err != nil {
		return plumbing.Hash{}, err
	}

	return head.Hash(), nil
}

// HardResetToCommit hard reset the worktree to the commit to clear changes
func (repo *Repository) HardResetToCommit(commit plumbing.Hash) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()
	worktree, err := repo.r.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Reset(&git.ResetOptions{
		Commit: commit,
		Mode:   git.HardReset,
	})

	logrus.Warn("repo is hard reset to head")

	return err
}

// Commit with username as slack caller name annotated by (Major Tom)
func (repo *Repository) Commit(filename, caller, message string) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()
	// TODO extract email and bot name as configuration
	return commit(repo, filename, fmt.Sprintf("%s(%s)", "Major Tom", caller), "mnews@mnews.tw", message)
}

func commit(repo *Repository, filename, name, email, message string) error {
	r := repo.r
	worktree, err := r.Worktree()
	if err != nil {
		return err
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

	logrus.Infof("commit message for %s:%s", repo.name, obj)
	return err
}

func (repo *Repository) Pull() error {
	repo.locker.Lock()
	defer repo.locker.Unlock()
	worktree, err := repo.r.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Pull(&git.PullOptions{
		Auth:          repo.authMethod,
		ReferenceName: plumbing.NewBranchReferenceName(repo.config.Branch),
		RemoteName:    "origin",
		SingleBranch:  true,
	})

	if err.Error() == "already up-to-date" {
		logrus.Infof("pulling repo, but it's already up-to-date")
		err = nil
	} else if err != nil {
		err = errors.Wrap(err, "pulling has error")
	}

	return err
}

func (repo *Repository) Push() error {
	repo.locker.Lock()
	defer repo.locker.Unlock()
	return repo.r.Push(&git.PushOptions{
		Auth: repo.authMethod,
	})
}

// FIXME
// GetRepository implicitly implies project equals repository. This is not a good practice, and it should be fixed in a good manner.
func GetRepository(project string, gitConfigs map[config.Repository]config.GitConfig) (r *Repository, err error) {

	gitConfig, isExisting := gitConfigs[config.Repository(project)]
	if !isExisting {
		return nil, errors.Errorf("the git config doesn't exist for the project(%s)", project)
	}

	// Get the singleton repository according to the project
	return getRepository(config.Repository(project), gitConfig)
}

func getRepository(project config.Repository, gitConfig config.GitConfig) (repo *Repository, err error) {

	switch project {
	// case "mm":
	// 	repo = mm
	case "tv":
		repo = tv
	// case "readr":
	// 	repo = readr
	default:
		return nil, errors.New("wrong project")
	}

	// Init git repo
	repo.once.Do(func() {
		repo.locker.Lock()
		defer repo.locker.Unlock()
		// Get the config according to the project
		repo.config = &gitConfig
		key, errRead := os.ReadFile(gitConfig.SSHKeyPath)
		if errRead != nil {
			err = errRead
			err = errors.Wrap(errRead, "reading ssh key failed")
			return
		}
		sshMethod, errSSH := ssh.NewPublicKeys(gitConfig.SSHKeyUser, key, "")
		if errSSH != nil {
			err = errors.Wrap(errSSH, "creating sshMethod from key failed")
			return
		}
		knownHostsFn, errKH := ssh.NewKnownHostsCallback(repo.config.SSHKnownhosts)
		if errKH != nil {
			err = errors.Wrap(errKH, "getting known_hosts file failed")
			return
		}
		sshMethod.HostKeyCallback = knownHostsFn
		repo.authMethod = sshMethod
		opt := git.CloneOptions{
			Auth:          repo.authMethod,
			ReferenceName: plumbing.NewBranchReferenceName(gitConfig.Branch),
			SingleBranch:  true,
			URL:           gitConfig.URL,
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
