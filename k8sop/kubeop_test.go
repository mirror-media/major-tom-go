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
