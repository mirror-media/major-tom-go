package command

import (
	"context"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	type args struct {
		ctx       context.Context
		textParts []string
	}
	tests := []struct {
		name        string
		args        args
		wantMessage []string
		wantErr     bool
	}{
		{
			name: "ask for wrong project",
			args: args{
				ctx:       context.TODO(),
				textParts: []string{"mn"},
			},
			// TODO change to config input
			wantMessage: []string{"call help"},
			wantErr:     true,
		},
		{
			name: "ask for projects",
			args: args{
				ctx:       context.TODO(),
				textParts: []string{},
			},
			// TODO change to config input
			wantMessage: []string{"The following projects are available: mm, tv, readr"},
		},
		{
			name: "ask for stages",
			args: args{
				ctx:       context.TODO(),
				textParts: []string{"mm"},
			},
			// TODO change to config input
			wantMessage: []string{"The following stages are available for mm: prod, staging, dev"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessage, err := List(tt.args.ctx, tt.args.textParts)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMessage, tt.wantMessage) {
				t.Errorf("List() = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}
