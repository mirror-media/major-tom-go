package gitop

import (
	"sync"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mirror-media/major-tom-go/config"
	"github.com/mirror-media/major-tom-go/internal/test"
)

func TestRepository_GetFile(t *testing.T) {
	repo, err := GetRepository("tv", test.GitConfigsTest)
	if err != nil {
		t.Error(err)
	}
	type fields struct {
		authMethod ssh.AuthMethod
		config     *config.GitConfig
		once       *sync.Once
		r          *git.Repository
		locker     *sync.Mutex
	}
	type args struct {
		filenamePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "dev for tv",
			fields: fields{
				authMethod: repo.authMethod,
				config: &config.GitConfig{
					Branch:        "test/majortom",
					SSHKeyPath:    "Users/chiu/dev/mtv/major-tom-go/configs/ssh/identity",
					SSHKeyUser:    "mnews@mnews.tw",
					SSHKnownhosts: "/Users/chiu/dev/mtv/major-tom-go/configs/ssh/known_hosts",
					URL:           "ssh://source.developers.google.com:2022/p/mirror-tv-275709/r/helm",
				},
				once:   repo.once,
				r:      repo.r,
				locker: repo.locker,
			},
			args: args{
				filenamePath: "cms/values-prod.yaml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				authMethod: tt.fields.authMethod,
				config:     tt.fields.config,
				once:       tt.fields.once,
				r:          tt.fields.r,
				locker:     tt.fields.locker,
			}
			_, err := repo.GetFile(tt.args.filenamePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
