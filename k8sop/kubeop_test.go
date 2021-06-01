package k8sop

import (
	"context"
	"reflect"
	"testing"
)

func TestGetResource(t *testing.T) {
	type args struct {
		ctx       context.Context
		namespace string
		name      string
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
				ctx:       context.TODO(),
				namespace: "default",
				name:      "www",
			},
			want: map[string]int{
				"Phase: Running; Ready: True": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResource(tt.args.ctx, tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
