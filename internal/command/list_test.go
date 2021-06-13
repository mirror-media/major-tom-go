package command

import (
	"context"
	"reflect"
	"testing"

	"github.com/mirror-media/major-tom-go/config"
)

func TestList(t *testing.T) {
	type args struct {
		ctx            context.Context
		clusterConfigs config.K8S
		textParts      []string
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
				ctx:            context.TODO(),
				clusterConfigs: clusterConfigs,
				textParts:      []string{"mn"},
			},
			wantMessage: []string{"call help"},
			wantErr:     true,
		},
		{
			name: "ask for projects",
			args: args{
				ctx:            context.TODO(),
				clusterConfigs: clusterConfigs,
				textParts:      []string{},
			},
			wantMessage: []string{"The following projects are available: mm, readr, tv"},
		},
		{
			name: "ask for stages",
			args: args{
				ctx:            context.TODO(),
				clusterConfigs: clusterConfigs,
				textParts:      []string{"mm"},
			},
			wantMessage: []string{"The following stages are available for mm: dev, prod, staging"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessage, err := List(tt.args.ctx, tt.args.clusterConfigs, tt.args.textParts)
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
