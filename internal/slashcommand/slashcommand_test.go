package slashcommand

import (
	"context"
	"reflect"
	"testing"

	"github.com/mirror-media/major-tom-go/v2/config"
)

// FIXME we need proper test path
var clusterConfigs = config.K8S{
	"mm": {
		"prod":    "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"staging": "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"dev":     "/Users/chiu/dev/mtv/major-tom-go/configs/config",
	},
	"tv": {
		"prod":    "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"staging": "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"dev":     "/Users/chiu/dev/mtv/major-tom-go/configs/config",
	},
	"readr": {
		"prod": "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"dev":  "/Users/chiu/dev/mtv/major-tom-go/configs/config",
	},
}

func TestRun(t *testing.T) {
	type args struct {
		ctx            context.Context
		clusterConfigs config.K8S
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
				clusterConfigs: clusterConfigs,
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
				clusterConfigs: clusterConfigs,
				cmd:            "/major-tom",
				txt:            "list",
			},
			wantMessages: []string{"The following projects are available: mm, tv, readr"},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, err := Run(tt.args.ctx, tt.args.clusterConfigs, tt.args.cmd, tt.args.txt)
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
