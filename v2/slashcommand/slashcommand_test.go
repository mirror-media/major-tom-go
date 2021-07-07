package slashcommand

import (
	"context"
	"reflect"
	"testing"

	"github.com/mirror-media/major-tom-go/v2/command"
	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/mirror-media/major-tom-go/v2/internal/test"
)

func TestRun(t *testing.T) {
	command.DeployWorker.Init(test.GitConfigsTest)
	type args struct {
		ctx            context.Context
		caller         string
		clusterConfigs config.K8S
		gitConfigs     map[config.Repository]config.GitConfig
		cmd            string
		txt            string
	}
	tests := []struct {
		name         string
		args         args
		wantMessages []string
		wantErr      bool
	}{

		{
			name: "wrong slashcommand",
			args: args{
				ctx:            context.TODO(),
				clusterConfigs: test.ConfigTest.ClusterConfigs,
				cmd:            "/mahjong-tom",
				txt:            "list",
			},
			// TODO change to helper function
			wantMessages: []string{"call help"},
			wantErr:      true,
		},
		{
			name: "list projects",
			args: args{
				ctx:            context.TODO(),
				clusterConfigs: test.ConfigTest.ClusterConfigs,
				cmd:            "/major-tom",
				txt:            "list",
			},
			wantMessages: []string{"The following projects are available: tv"},
			wantErr:      false,
		},
		{
			name: "deploy cms",
			args: args{
				ctx:            context.TODO(),
				caller:         "tester",
				clusterConfigs: test.ConfigTest.ClusterConfigs,
				cmd:            "/major-tom",
				txt:            "deploy tv prod cms pods:2 image:imageTag maxPods:4 minPods:1 autoScaling:true",
			},
			wantMessages: []string{"release(cms/prod): released by @tester", "", "Set autoscaling.enabled to true", "Set autoscaling.maxReplicas to 4", "Set autoscaling.minReplicas to 1", "Set image.tag to imageTag", "Set replicacount to 2", ""},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, err := Run(tt.args.ctx, tt.args.clusterConfigs, tt.args.cmd, tt.args.txt, tt.args.caller)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMessages, tt.wantMessages) {
				t.Errorf("Run() = %v, want %v", gotMessages, tt.wantMessages)
			}
		})
	}
}
