package types

import "testing"

func TestCommand_MatchesCommand(t *testing.T) {
	type fields struct {
		Command string
		Aliases []string
	}
	type args struct {
		command string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "MatchesCommand",
			fields: fields{
				Command: "test",
				Aliases: []string{"test2"},
			},
			args: args{
				command: "test",
			},
			want: true,
		},
		{
			name: "MatchesCommandNotFound",
			fields: fields{
				Command: "test",
				Aliases: []string{"3test2"},
			},
			args: args{
				command: "test2",
			},
			want: false,
		},
		{
			name: "MatchesCommandAlias",
			fields: fields{
				Command: "test",
				Aliases: []string{"test2"},
			},
			args: args{
				command: "test2",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Command{
				Command: tt.fields.Command,
				Aliases: tt.fields.Aliases,
			}
			if got := c.MatchesCommand(tt.args.command); got != tt.want {
				t.Errorf("Command.MatchesCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
