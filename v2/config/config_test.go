package config

import "testing"

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
