package assets

import (
	"io/fs"

	embeddedassets "github.com/homekit/homekit-cli/assets"
)

// Embedded exposes the build-time embedded filesystem.
func Embedded() fs.FS {
	return embeddedassets.FS
}

type AssetNamespace string

const (
	AssetNamespaceScripts    AssetNamespace = "scripts"
	AssetNamespaceTemplates  AssetNamespace = "templates"
	AssetNamespaceWorkspaces AssetNamespace = "workspaces"
)

func (a AssetNamespace) String() string {
	return string(a)
}
