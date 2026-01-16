package static

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var staticFiles embed.FS

// GetStaticFS returns the embedded static filesystem rooted at dist/
func GetStaticFS() (fs.FS, error) {
	return fs.Sub(staticFiles, "dist")
}
