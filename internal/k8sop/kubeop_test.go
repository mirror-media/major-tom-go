package k8sop

import (
	"context"
	"reflect"
	"testing"
)

func TestGetResource(t *testing.T) {
	type args struct {
		ctx            context.Context
		kubeConfigPath string
		namespace      string
		name           string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{
			name: "get www pods",
			args: args{
				ctx: context.TODO(),
				// FIXME we need proper path
				kubeConfigPath: "/Users/chiu/dev/mtv/major-tom-go/configs/config",
				namespace:      "default",
				name:           "www",
			},
			want: map[string]int{
				"master__272, Phase: Running, Ready: True": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPodInfo(tt.args.ctx, tt.args.kubeConfigPath, tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPodInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDeploymentInfo(t *testing.T) {
	type args struct {
		ctx            context.Context
		kubeConfigPath string
		namespace      string
		name           string
	}
	tests := []struct {
		name    string
		args    args
		want    DeploymentInfo
		wantErr bool
	}{
		// FIXME we need proper test cases
		{
			name: "get www deployment info",
			args: args{
				ctx:            context.Background(),
				kubeConfigPath: "/Users/chiu/dev/mtv/major-tom-go/configs/config",
				namespace:      "default",
				name:           "www",
			},
			want: DeploymentInfo{
				Available: 1,
				ImageTag:  "master__272",
				Ready:     1,
				Updated:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDeploymentInfo(tt.args.ctx, tt.args.kubeConfigPath, tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDeploymentInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDeploymentInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listServices(t *testing.T) {
	type args struct {
		ctx            context.Context
		kubeConfigPath string
		namespace      string
	}
	tests := []struct {
		name            string
		args            args
		wantReleaseInfo []ReleaseInfo
		wantErr         bool
	}{
		// FIXME we need proper test cases
		{
			name: "list dev",
			args: args{
				ctx:            context.TODO(),
				kubeConfigPath: "/Users/chiu/dev/mtv/major-tom-go/configs/config",
			},
			wantReleaseInfo: []ReleaseInfo{
				{Status: "deployed", Name: "cms", Namespace: "default"},
				{Status: "deployed", Name: "cronjobs", Namespace: "cron"},
				{Status: "failed", Name: "elasticsearch", Namespace: "textsearch"},
				{Status: "deployed", Name: "graphql-external", Namespace: "default"},
				{Status: "deployed", Name: "graphql-internal", Namespace: "default"},
				{Status: "deployed", Name: "nginx-web", Namespace: "default"},
				{Status: "deployed", Name: "redis-cluster", Namespace: "redis"},
				{Status: "deployed", Name: "secrets", Namespace: "flux"},
				{Status: "deployed", Name: "www", Namespace: "default"},
				{Status: "deployed", Name: "yt-relay", Namespace: "default"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReleaseInfo, err := listReleases(tt.args.ctx, tt.args.kubeConfigPath, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("listReleases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReleaseInfo, tt.wantReleaseInfo) {
				t.Errorf("listReleases() = %v, want %v", gotReleaseInfo, tt.wantReleaseInfo)
			}
		})
	}
}

func TestGetDeploymentInfo(t *testing.T) {
	type args struct {
		ctx            context.Context
		kubeConfigPath string
		name           string
	}
	tests := []struct {
		name    string
		args    args
		want    DeploymentInfo
		wantErr bool
	}{
		// FIXME we need proper test cases
		{
			name: "wip",
			args: args{
				ctx:            context.TODO(),
				kubeConfigPath: "/Users/chiu/dev/mtv/major-tom-go/configs/config",
				name:           "yt-relay",
			},
			want: DeploymentInfo{
				Available: 1,
				ImageTag:  "master__59",
				Ready:     1,
				Updated:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDeploymentInfo(tt.args.ctx, tt.args.kubeConfigPath, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeploymentInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDeploymentInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
