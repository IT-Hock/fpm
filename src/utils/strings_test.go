package utils

import "testing"

func TestHashString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty string",
			args: args{s: ""},
			want: "2166136261",
		},
		{
			name: "string with length 1",
			args: args{s: "a"},
			want: "3826002220",
		},
		{
			name: "string with length 2",
			args: args{s: "ab"},
			want: "1294271946",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashString(tt.args.s); got != tt.want {
				t.Errorf("HashString() = %v, want %v", got, tt.want)
			}
		})
	}
}
