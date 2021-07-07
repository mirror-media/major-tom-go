package command

import (
	"context"
	"reflect"
	"testing"

	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/mirror-media/major-tom-go/v2/internal/test"
)

func TestInfo(t *testing.T) {
	type args struct {
		ctx            context.Context
		clusterConfigs config.K8S
		textParts      []string
	}
	tests := []struct {
		name         string
		args         args
		wantMessages []string
		wantErr      bool
	}{
		{
			name: "ask for wrong project",
			args: args{
				ctx:            context.TODO(),
				clusterConfigs: test.ConfigTest.ClusterConfigs,
				textParts:      []string{"mn"},
			},
			wantMessages: []string{"call help"},
			wantErr:      true,
		},
		{
			name: "ask for wrong service",
			args: args{
				ctx:            context.TODO(),
				clusterConfigs: test.ConfigTest.ClusterConfigs,
				textParts:      []string{"tv", "prod", "TaiwanNumberOne"},
			},
			wantMessages: []string{"call list"},
			wantErr:      true,
		},
		{
			name: "ask for yt-relay",
			args: args{
				ctx:            context.TODO(),
				clusterConfigs: test.ConfigTest.ClusterConfigs,
				textParts:      []string{"tv", "prod", "yt-relay"},
			},
			wantMessages: []string{"yt-relay\n\tImageTag: master__59\n\tAvailable pods: 1\n\tReady pods: 1\n\tUpdated pods: 1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, err := Info(tt.args.ctx, tt.args.clusterConfigs, tt.args.textParts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMessages, tt.wantMessages) {
				t.Errorf("Info() = %v, want %v", gotMessages, tt.wantMessages)
			}
		})
	}
}
