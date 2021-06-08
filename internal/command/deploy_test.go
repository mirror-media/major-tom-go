package command

import (
	"context"
	"reflect"
	"testing"

	"github.com/mirror-media/major-tom-go/v2/config"
	mjcontext "github.com/mirror-media/major-tom-go/v2/internal/context"
)

func Test_deploy(t *testing.T) {
	ch := make(chan CommandResponse, 100)
	ctx := context.WithValue(context.TODO(), mjcontext.ResponseChannel, ch)
	type args struct {
		ctx            context.Context
		clusterConfigs config.K8S
		textParts      []string
		caller         string
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
				caller:         "@tester",
				clusterConfigs: clusterConfigs,
				ctx:            ctx,
			},
			wantMessages: []string{"call help"},
			wantErr:      true,
		},
		{
			name: "dev",
			args: args{
				caller:         "@tester",
				clusterConfigs: clusterConfigs,
				ctx:            ctx,
				textParts:      []string{"tv", "prod", "yt-relay", "image:11", "pods:23"},
			},
			wantMessages: []string{"release(yt-relay/prod): released by @tester", "", "Set image.tag to 11", "Set replicacount to 23", ""},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deploy(tt.args.ctx, tt.args.clusterConfigs, tt.args.textParts, tt.args.caller)
			ch := tt.args.ctx.Value(mjcontext.ResponseChannel).(chan CommandResponse)
			gotResponse := <-ch
			err := gotResponse.Error
			if (err != nil) != tt.wantErr {
				t.Errorf("Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotMessages := gotResponse.Messages
			if !reflect.DeepEqual(gotMessages, tt.wantMessages) {
				t.Errorf("Info() = %v, want %v", gotMessages, tt.wantMessages)
			}
		})
	}
}
