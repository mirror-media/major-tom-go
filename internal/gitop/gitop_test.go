package gitop

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
)

func TestRepository_GetFile(t *testing.T) {
	repo, err := GetRepository("tv")
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	type fields struct {
		config map[string]string
		once   *sync.Once
		r      *git.Repository
	}
	type args struct {
		filenamePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    billy.File
		wantErr bool
	}{
		{
			name: "dev for tv",
			fields: fields{
				config: gitConfig["tv"],
				once:   repo.once,
				r:      repo.r,
			},
			args: args{
				filenamePath: "cms/values-prod.yaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				config: tt.fields.config,
				once:   tt.fields.once,
				r:      tt.fields.r,
			}
			got, err := repo.GetFile(tt.args.filenamePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
