package types

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"regexp"
	"testing"
	"time"
)

//go:embed packagesMap.json
var packagesMap []byte

func Benchmark_FindPackages(b *testing.B) {
	var p Packages
	err := json.Unmarshal(packagesMap, &p)
	if err != nil {
		return
	}

	expression := regexp.MustCompile("not-found")

	elapsed := 0
	for i := 0; i < b.N; i++ {
		start := time.Now()
		_, _ = p.FindPackages(expression)
		elapsed += int(time.Since(start).Nanoseconds())
	}
	b.Logf("Iterations: %d, Average time: %s, Total time: %s", b.N, formatTime(int64(elapsed/b.N)), formatTime(int64(elapsed)))
	b.Logf("\n")

	elapsed = 0
	for i := 0; i < b.N; i++ {
		start := time.Now()
		_, _ = p.FindPackagesFast("not-found")
		elapsed += int(time.Since(start).Nanoseconds())
	}
	b.Logf("Iterations: %d, Average time: %s, Total time: %s", b.N, formatTime(int64(elapsed/b.N)), formatTime(int64(elapsed)))
	b.Logf("\n")
}

func formatTime(ns int64) string {
	return time.Duration(ns).String()
}

func TestPackages_FindPackage(t *testing.T) {
	type fields struct {
		Packages map[string]Package
		Themes   map[string]Package
	}
	type args struct {
		name string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    *Package
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packages{
				Packages: tt.fields.Packages,
				Themes:   tt.fields.Themes,
			}
			got, err := p.FindPackage(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPackage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackages_FindPackages(t *testing.T) {
	type fields struct {
		Packages map[string]Package
		Themes   map[string]Package
	}
	type args struct {
		expression *regexp.Regexp
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    []Package
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packages{
				Packages: tt.fields.Packages,
				Themes:   tt.fields.Themes,
			}
			got, err := p.FindPackages(tt.args.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPackages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPackages() got = %v, want %v", got, tt.want)
			}
		})
	}
}
