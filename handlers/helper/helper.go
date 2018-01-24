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

var indexFileNames = []string{"index.html", "index.htm"}

func (d FSDir) Open(name string) (http.File, error) {
	logs.Debugf("The name is %s", name)
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) ||
		strings.Contains(name, "\x00") {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	realPath := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	logs.Debugf("The real path is %s", realPath)
	fileStat, err := os.Stat(realPath)
	if err != nil {
		if os.IsNotExist(err) {
			logs.Errorf("File %s does not exists", realPath)
			return nil, err
		}
		logs.Errorf("Failed to get status for file %s, the error is %#v", realPath, err)
		return nil, err
	}
	// Check index file
	if fileStat.IsDir() {
		indexFilePath, err := getIndexFilePath(realPath)
		if err != nil {
			logs.Errorf("Failed to get index file path from dir %s, the error is %#v", realPath, err)
			return nil, err
		}
		if name == "/" {
			return getFileFromPath(realPath)
		}
		logs.Debugf("The index file path is %s", indexFilePath)
		return getFileFromPath(indexFilePath)
	}

	if !fileStat.IsDir() {
		return getFileFromPath(realPath)
	}
	logs.Errorf("The path %s is dir, which is not allowed", realPath)
	return nil, os.ErrPermission
}

func getIndexFilePath(dirPath string) (string, error) {
	logs.Debugf("Get index file path from dir %s", dirPath)
	for _, indexFileName := range indexFileNames {
		indexFilePath := filepath.Join(dirPath, indexFileName)
		logs.Debugf("Check index file %s", indexFilePath)
		osStat, err := os.Stat(indexFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		if !osStat.IsDir() {
			return indexFilePath, nil
		}
	}
	return "", os.ErrNotExist
}

func getFileFromPath(realPath string) (http.File, error) {
	f, err := os.Open(realPath)
	//defer f.Close()
	if err != nil {
		logs.Errorf("Failed to open file %s, the error is %#v", realPath, err)
		return nil, err
	}
	return f, nil
}
