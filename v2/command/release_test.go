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

func TestRelease(t *testing.T) {

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
			name: "no texts",
			args: args{
				caller:  "+tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "invalid repo",
			args: args{
				caller:  "+tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"invalidRepo", "project=tv", "image=" + newTag},
			},
			wantErr: true,
		},
		{
			name: "invalid stage",
			args: args{
				caller:  "+tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"openwarehouse", "env=dev", "image=" + newTag},
				message: "deploy openwarehouse env=dev image-tag=" + newTag,
			},
			wantErr: true,
		},
		{
			name: "invalid project",
			args: args{
				caller:  "+tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"openwarehouse", "project=invalid", "image-tag=" + newTag},
				message: "release openwarehouse project=invalid image-tag=" + newTag,
			},
			wantErr: true,
		},
		{
			name: "release openwarehouse tv image-tag",
			args: args{
				caller:  "+tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"openwarehouse", "project=tv", "image-tag=" + newTag},
				message: "release openwarehouse project=tv image-tag=" + newTag,
			},
			wantMessages: []string{"deploy(openwarehouse/prod/tv): deployed by +tester", "", "Set image-tag(images.0.newTag) to " + newTag, "by \"release openwarehouse project=tv image-tag=" + newTag + "\""},
		},
		{
			name: "release mirror-tv-nuxt tv image-tag",
			args: args{
				caller:  "+tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"mirror-tv-nuxt", "project=tv", "image-tag=" + newTag},
				message: "release mirror-tv-nuxt project=tv image-tag=" + newTag,
			},
			wantMessages: []string{"deploy(mirror-tv-nuxt/prod/tv): deployed by +tester", "", "Set image-tag(images.0.newTag) to " + newTag, "by \"release mirror-tv-nuxt project=tv image-tag=" + newTag + "\""},
		},
		{
			name: "can't release without project",
			args: args{
				caller:  "+tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"openwarehouse", "env=prod", "image-tag=" + newTag},
				message: "release openwarehouse env=prod image-tag=" + newTag,
			},
			wantErr: true,
		},
		{
			name: "can't release with any unsupported extra argument",
			args: args{
				caller:  "+tester",
				k8sRepo: test.K8sRepo,
				ctx:     context.TODO(),
				texts:   []string{"openwarehouse", "project=tv", "env=prod", "image-tag=" + newTag},
				message: "release openwarehouse env=prod image-tag=" + newTag,
			},
			wantErr: true,
		},
	}
	DeployWorker.Set(test.K8sRepo.GitConfig)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, err := Release(tt.args.ctx, tt.args.k8sRepo, tt.args.texts, tt.args.message, tt.args.caller)
			if (err != nil) != tt.wantErr {
				t.Errorf("Release() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMessages, tt.wantMessages) {
				t.Errorf("Release() = %v, want %v", gotMessages, tt.wantMessages)
			}
		})
	}
}
