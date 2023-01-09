package types

import (
	"testing"
)

func TestPackage_fetchFromGitHub(t *testing.T) {
	type fields struct {
		Name         string
		Description  string
		Image        string
		Tags         []string
		Author       string
		Repository   string
		Dependencies []string
		IsTheme      bool
	}
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "fetchFromGitHub",
			fields: fields{
				Name:        "test",
				Description: "test",
				Image:       "test",
				Tags:        []string{"test"},
				Author:      "test",
				Repository:  "https://github.com/IT-Hock/fpm",
			},
			args: args{
				version: "v0.0.1",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Package{
				Name:         tt.fields.Name,
				Description:  tt.fields.Description,
				Image:        tt.fields.Image,
				Tags:         tt.fields.Tags,
				Author:       tt.fields.Author,
				Repository:   tt.fields.Repository,
				Dependencies: tt.fields.Dependencies,
				IsTheme:      tt.fields.IsTheme,
			}
			got, got1, err := p.fetchFromGitHub(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchFromGitHub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetchFromGitHub() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("fetchFromGitHub() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPackage_fetchFromGitLab(t *testing.T) {
	type fields struct {
		Name         string
		Description  string
		Image        string
		Tags         []string
		Author       string
		Repository   string
		Dependencies []string
		IsTheme      bool
	}
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "fetchFromGitLab",
			fields: fields{
				Name:        "test",
				Description: "test",
				Image:       "test",
				Tags:        []string{"test"},
				Author:      "test",
				Repository:  "https://gitlab.com/IT-Hock/fpm",
			},
			args: args{
				version: "v0.0.1",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Package{
				Name:         tt.fields.Name,
				Description:  tt.fields.Description,
				Image:        tt.fields.Image,
				Tags:         tt.fields.Tags,
				Author:       tt.fields.Author,
				Repository:   tt.fields.Repository,
				Dependencies: tt.fields.Dependencies,
				IsTheme:      tt.fields.IsTheme,
			}
			got, got1, err := p.fetchFromGitLab(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchFromGitLab() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetchFromGitLab() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("fetchFromGitLab() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
