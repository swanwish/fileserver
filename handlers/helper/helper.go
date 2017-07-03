package helper

import (
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/swanwish/go-common/logs"
)

type FSDir string

func (d FSDir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) ||
		strings.Contains(name, "\x00") {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	realPath := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	if fileStat, err := os.Stat(realPath); err == nil && !fileStat.IsDir() {
		f, err := os.Open(realPath)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	logs.Errorf("The path %s is dir, which is not allowed", realPath)
	return nil, os.ErrPermission
}
