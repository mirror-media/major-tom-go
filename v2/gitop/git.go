package gitop

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
)

func cloneGitRepo(opt git.CloneOptions) (*git.Repository, error) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &opt)

	if err != nil {
		return nil, err
	}

	// Gets the HEAD history from HEAD, just like this command:
	log.Info("git log")

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	// Gets the HEAD history from HEAD, just like this command:
	log.Infof("%s is at head %s", ref.Name(), ref.Hash())
	return r, nil
}
