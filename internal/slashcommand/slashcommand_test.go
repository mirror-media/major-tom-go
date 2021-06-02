package slashcommand

import (
	"context"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		ctx context.Context
		cmd string
		txt string
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
				ctx: context.TODO(),
				cmd: "/mahjong-tom",
				txt: "list",
			},
			// TODO change to helper function
			wantMessages: []string{"call help"},
			wantErr:      true,
		},
		{
			name: "list projects",
			args: args{
				ctx: context.TODO(),
				cmd: "/major-tom",
				txt: "list",
			},
			// TODO change to config input
			wantMessages: []string{"The following projects are available: mm, tv, readr"},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, err := Run(tt.args.ctx, tt.args.cmd, tt.args.txt)
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
