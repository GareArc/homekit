package assets

import (
	"io/fs"

	embeddedassets "github.com/homekit/homekit-cli/assets"
)

// Embedded exposes the build-time embedded filesystem.
func Embedded() fs.FS {
	return embeddedassets.FS
}
