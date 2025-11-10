package assets

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Manager provides accessors for embedded and user-supplied assets.
type Manager struct {
	embedded fs.FS
	override string
}

// NewManager constructs a new asset manager.
func NewManager(embedded fs.FS, overrideDir string) *Manager {
	return &Manager{
		embedded: embedded,
		override: overrideDir,
	}
}

// List returns sorted asset names under the provided namespace (scripts/templates).
func (m *Manager) List(namespace string) ([]string, error) {
	base := strings.Trim(namespace, "/")

	set := map[string]struct{}{}

	if err := fs.WalkDir(m.embedded, base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel := strings.TrimPrefix(path, base+"/")
		set[rel] = struct{}{}
		return nil
	}); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if m.override != "" {
		root := filepath.Join(m.override, base)
		_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			set[filepath.ToSlash(rel)] = struct{}{}
			return nil
		})
	}

	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

// Open returns a read handle for an asset, preferring overrides.
func (m *Manager) Open(namespace, name string) (fs.File, error) {
	if m.override != "" {
		path := filepath.Join(m.override, namespace, name)
		if f, err := os.Open(path); err == nil {
			return f, nil
		}
	}
	return m.embedded.Open(filepath.ToSlash(filepath.Join(namespace, name)))
}

// OpenBytes returns the content of an asset as a byte slice.
func (m *Manager) OpenBytes(namespace AssetNamespace, name string) ([]byte, error) {
	src, err := m.Open(namespace.String(), name)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	return io.ReadAll(src)
}

// Export copies the asset to the destination directory.
func (m *Manager) Export(namespace, name, destDir string) (string, error) {
	src, err := m.Open(namespace, name)
	if err != nil {
		return "", err
	}
	defer src.Close()

	target := filepath.Join(destDir, name)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", err
	}

	out, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}
	if err := out.Chmod(0o755); err != nil {
		return "", err
	}
	return target, nil
}

// Verify calculates a checksum for an asset.
func (m *Manager) Verify(namespace, name string) (string, error) {
	file, err := m.Open(namespace, name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Filesystem exposes the embedded filesystem for direct access.
func (m *Manager) Filesystem() fs.FS {
	return m.embedded
}
