package command

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/mirror-media/major-tom-go/v2/internal/test"
)

func TestDeploy(t *testing.T) {

	randomBytes := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	_, _ = rand.Read(randomBytes)
	h := md5.New()
	h.Write(randomBytes)
	newTag := "test-new-tag-" + fmt.Sprintf("%x", h.Sum(nil))[:5]

	type args struct {
		ctx     context.Context
		k8sRepo config.KubernetesConfigsRepo
		texts   []string
		message string
		caller  string
	}

	tests := []struct {
		name         string
		args         args
		wantMessages []string
		wantErr      bool
	}{
		{
			name: "no textParts",
			args: args{
				caller:  "@tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "invalidRepo",
			args: args{
				caller:  "@tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"invalidRepo", "env=dev", "image=" + newTag},
			},
			wantErr: true,
		},
		{
			name: "deploy env openwarehouse image-tag",
			args: args{
				caller:  "@tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"openwarehouse", "env=dev", "image-tag=" + newTag},
				message: "deploy openwarehouse env=dev image-tag=" + newTag,

			},
			wantMessages: []string{"deploy(openwarehouse/dev): deployed by @tester", "", "Set image-tag(images.0.newTag) to " + newTag, "by \"deploy openwarehouse env=dev image-tag=" + newTag + "\""},
		},
		{
			name: "deploy env mirror-tv-nuxt image-tag",
			args: args{
				caller:  "@tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"mirror-tv-nuxt", "env=dev", "image-tag=" + newTag},
				message: "deploy mirror-tv-nuxt env=dev image-tag=" + newTag,
			},
			wantMessages: []string{"deploy(mirror-tv-nuxt/dev): deployed by @tester", "", "Set image-tag(images.0.newTag) to " + newTag, "by \"deploy mirror-tv-nuxt env=dev image-tag=" + newTag + "\""},
		},
		{
			name: "can't deploy prod image-tag",
			args: args{
				caller:  "@tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"openwarehouse", "env=prod", "image-tag=" + newTag},
				message: "deploy openwarehouse env=prod image-tag=" + newTag,
			},
			wantErr: true,
		},
	}
	DeployWorker.Set(test.K8sRepo.GitConfig)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, err := Deploy(tt.args.ctx, tt.args.k8sRepo, tt.args.texts, tt.args.message, tt.args.caller)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deploy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMessages, tt.wantMessages) {
				for _, item := range gotMessages {
					fmt.Println("'" + item + "'")
				}
				for _, item := range tt.wantMessages {
					fmt.Println("'" + item + "'")
				}
				t.Errorf("Deploy() = %+v, want %+v", gotMessages, tt.wantMessages)
			}
		})
	}
}
