package pathformat

import (
	"os"
	"path/filepath"
)

func RenderFullPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(Pwd(), path)
}

func Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func Pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return pwd
}

func Join(elem ...string) string {
	return filepath.Join(elem...)
}

func Clean(path string) string {
	return filepath.Clean(path)
}

func MakeDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func Base(path string) string {
	return filepath.Base(path)
}
