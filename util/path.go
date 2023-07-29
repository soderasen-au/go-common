package util

import (
	"io/ioutil"
	"os"
	"path"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ListSubFolders(path string) ([]string, error) {
	exists, err := Exists(path)
	if err != nil || !exists {
		return nil, err
	}

	subfolders := make([]string, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			subfolders = append(subfolders, f.Name())
		}
	}
	return subfolders, nil
}

// will return full path.
func ListFiles(folder string) ([]string, *Result) {
	exists, err := Exists(folder)
	if err != nil || !exists {
		return nil, Error("Exists", err)
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, Error("ReadDir", err)
	}

	ret := make([]string, 0)
	for _, file := range files {
		if file.Type().IsRegular() {
			ret = append(ret, path.Join(folder, file.Name()))
		}
	}

	return ret, nil
}

func MaybeCreate(folder string) error {
	return os.MkdirAll(folder, 0777)
}
