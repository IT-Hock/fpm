package resources

import (
	"embed"
	_ "embed"
)

//go:embed files
var fishFiles embed.FS
