package config

import (
	"reflect"
	"testing"
)

// The config test is important because they determine the files to be changed for different operations
// They are the core of the major tom

func Test_contains(t *testing.T) {
	type args struct {
		s      []string
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil slice return false",
			args: args{
				s:      nil,
				target: "target",
			},
			want: false,
		},
		{
			name: "empty target but not match",
			args: args{
				s:      []string{"no roots", "Alice Merton"},
				target: "",
			},
			want: false,
		},
		{
			name: "partial match return false",
			args: args{
				s:      []string{"no roots", "Alice Merton"},
				target: "Alice",
			},
			want: false,
		},
		{
			name: "match return true",
			args: args{
				s:      []string{"no roots", "Alice Merton"},
				target: "Alice Merton",
			},
			want: true,
		},
		{
			name: "it's case sensitive",
			args: args{
				s:      []string{"no roots", "Alice Merton"},
				target: "alice Merton",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.s, tt.args.target); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodebase_GetServices(t *testing.T) {
	type fields struct {
		Projects []string
		Repo     string
		Services []string
		Stages   []string
		Type     int8
	}
	tests := []struct {
		name         string
		fields       fields
		wantServices []Service
		wantErr      bool
	}{
		{
			name: "type 2",
			fields: fields{
				Type:     2,
				Repo:     "repoXYZ",
				Projects: []string{"p1", "p2"},
				Services: []string{"s1", "s2"},
			},
			wantServices: []Service{
				{
					Name:          "repoXYZ-p1-s1",
					Repo:          "repoXYZ",
					SimpleService: "s1",
				},
				{
					Name:          "repoXYZ-p1-s2",
					Repo:          "repoXYZ",
					SimpleService: "s2",
				},
				{
					Name:          "repoXYZ-p2-s1",
					Repo:          "repoXYZ",
					SimpleService: "s1",
				},
				{
					Name:          "repoXYZ-p2-s2",
					Repo:          "repoXYZ",
					SimpleService: "s2",
				},
			},
		},
		{
			name: "type 2 answer is sorted by name",
			fields: fields{
				Type:     2,
				Repo:     "repoXYZ",
				Projects: []string{"p2", "p1"},
				Services: []string{"s2", "s1"},
			},
			wantServices: []Service{
				{
					Name:          "repoXYZ-p1-s1",
					Repo:          "repoXYZ",
					SimpleService: "s1",
				},
				{
					Name:          "repoXYZ-p1-s2",
					Repo:          "repoXYZ",
					SimpleService: "s2",
				},
				{
					Name:          "repoXYZ-p2-s1",
					Repo:          "repoXYZ",
					SimpleService: "s1",
				},
				{
					Name:          "repoXYZ-p2-s2",
					Repo:          "repoXYZ",
					SimpleService: "s2",
				},
			},
		},
		{
			name: "stages of repo doesn't change the service",
			fields: fields{
				Type:     2,
				Repo:     "repoXYZ",
				Stages:   []string{"dev", "staging", "prod"},
				Projects: []string{"p2", "p1"},
				Services: []string{"s2", "s1"},
			},
			wantServices: []Service{
				{
					Name:          "repoXYZ-p1-s1",
					Repo:          "repoXYZ",
					SimpleService: "s1",
				},
				{
					Name:          "repoXYZ-p1-s2",
					Repo:          "repoXYZ",
					SimpleService: "s2",
				},
				{
					Name:          "repoXYZ-p2-s1",
					Repo:          "repoXYZ",
					SimpleService: "s1",
				},
				{
					Name:          "repoXYZ-p2-s2",
					Repo:          "repoXYZ",
					SimpleService: "s2",
				},
			},
		},
		{
			name: "type 1 doesn't have simple service name",
			fields: fields{
				Type:   1,
				Repo:   "repoXYZ",
				Stages: []string{"dev", "staging", "prod"},
			},
			wantServices: []Service{
				{
					Name: "repoXYZ",
					Repo: "repoXYZ",
				},
			},
		},
		{
			name: "stages doesn't change service for type 1 either",
			fields: fields{
				Type: 1,
				Repo: "repoXYZ",
			},
			wantServices: []Service{
				{
					Name: "repoXYZ",
					Repo: "repoXYZ",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Codebase{
				Projects: tt.fields.Projects,
				Repo:     tt.fields.Repo,
				Services: tt.fields.Services,
				Stages:   tt.fields.Stages,
				Type:     tt.fields.Type,
			}
			gotServices, err := c.GetServices()
			if (err != nil) != tt.wantErr {
				t.Errorf("Codebase.GetServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotServices, tt.wantServices) {
				t.Errorf("Codebase.GetServices() = %v, want %v", gotServices, tt.wantServices)
			}
		})
	}
}

func TestCodebase_getType1RepoPath(t *testing.T) {
	type fields struct {
		Projects []string
		Repo     string
		Services []string
		Stages   []string
		Type     int8
	}
	type args struct {
		filename string
		stage    string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantPath string
		wantErr  bool
	}{
		{
			name: "path for type 2 will report error",
			fields: fields{
				Type:     2,
				Repo:     "repoXYZ",
				Stages:   []string{"ss1", "ss2"},
				Projects: []string{"p1", "p2"},
				Services: []string{"s1", "s2"},
			},
			args: args{
				filename: "filename.ext",
				stage:    "ss1",
			},
			wantPath: "repoXYZ/overlays/ss1/filename.ext",
			wantErr:  true,
		},
		{
			name: "path for type 1",
			fields: fields{
				Type:   1,
				Repo:   "repoXYZ",
				Stages: []string{"ss1", "ss2"},
			},
			args: args{
				filename: "filename.ext",
				stage:    "ss1",
			},
			wantPath: "repoXYZ/overlays/ss1/filename.ext",
		},
		{
			name: "wrong stage will still return a path for type 1",
			fields: fields{
				Type:   1,
				Repo:   "repoXYZ",
				Stages: []string{"ss1", "ss2"},
			},
			args: args{
				filename: "filename.ext",
				stage:    "s1",
			},
			wantPath: "repoXYZ/overlays/s1/filename.ext",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Codebase{
				Projects: tt.fields.Projects,
				Repo:     tt.fields.Repo,
				Services: tt.fields.Services,
				Stages:   tt.fields.Stages,
				Type:     tt.fields.Type,
			}
			gotPath, err := c.getType1RepoPath(tt.args.filename, tt.args.stage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Codebase.getType1RepoPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("Codebase.getType1RepoPath() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
