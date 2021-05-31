package k8sop

import (
	"reflect"
	"testing"
)

func TestListHelmReleaseTV(t *testing.T) {
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
				"tag": "master__272",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHelmReleaseImageTag(tt.args.name)
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
