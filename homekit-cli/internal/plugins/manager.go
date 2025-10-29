package plugins

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Manager discovers executable plugins matching a naming convention.
type Manager struct {
	Prefix      string
	SearchPaths []string
}

// NewManager creates a plugin manager with defaults.
func NewManager(prefix string, searchPaths []string) *Manager {
	if len(searchPaths) == 0 {
		searchPaths = filepath.SplitList(os.Getenv("PATH"))
	}
	return &Manager{Prefix: prefix, SearchPaths: searchPaths}
}

// Descriptor describes an installed plugin.
type Descriptor struct {
	Name string
	Path string
}

// Discover returns all visible plugins.
func (m *Manager) Discover() ([]Descriptor, error) {
	if m.Prefix == "" {
		return nil, errors.New("plugin prefix must be set")
	}

	seen := map[string]Descriptor{}
	for _, root := range m.SearchPaths {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !strings.HasPrefix(name, m.Prefix) {
				continue
			}
			full := filepath.Join(root, name)
			if !isExecutable(entry, full) {
				continue
			}
			short := strings.TrimPrefix(name, m.Prefix)
			short = strings.TrimPrefix(short, "-")
			seen[short] = Descriptor{
				Name: short,
				Path: full,
			}
		}
	}

	result := make([]Descriptor, 0, len(seen))
	for _, d := range seen {
		result = append(result, d)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

// ExecProxy constructs an *exec.Cmd for the plugin according to the provided args.
func (m *Manager) ExecProxy(descriptor Descriptor, args []string, env []string) *exec.Cmd {
	cmd := exec.Command(descriptor.Path, args...)
	cmd.Env = append(os.Environ(), env...)
	return cmd
}

func isExecutable(entry fs.DirEntry, fullPath string) bool {
	info, err := entry.Info()
	if err != nil {
		return false
	}
	mode := info.Mode()
	return mode.IsRegular() && mode&0o111 != 0
}
