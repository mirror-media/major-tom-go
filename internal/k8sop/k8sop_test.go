package k8sop

import (
	"reflect"
	"testing"
)

func TestGetHelmReleaseInfo(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "www release",
			args: args{
				name: "www",
			},
			want: map[string]string{
				// FIXME it's a temp test case
				DeploymentTagKey:    "master__272",
				DeploymentStatusKey: "deployed",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHelmReleaseInfo(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListHelmReleaseTV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListHelmReleaseTV() = %v, want %v", got, tt.want)
			}
		})
	}
}
