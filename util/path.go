package util

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	subfolders := make([]string, 0)
	files, err := os.ReadDir(path)
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
	if err != nil {
		return nil, Error("Exists", err)
	}
	if !exists {
		return nil, MsgError("Exists", "folder does not exist: "+folder)
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

func FilterFiles(folder, pattern string) ([]string, *Result) {
	ret := make([]string, 0)
	err := filepath.WalkDir(folder, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't parse file path %s: %s", file, err.Error())
		}
		if entry.Type().IsDir() {
			return nil
		}

		if entry.Type().IsRegular() {
			if matched, err := filepath.Match(pattern, filepath.Base(file)); err != nil {
				return fmt.Errorf("can't match file path %s: %s", file, err.Error())
			} else if matched {
				ret = append(ret, file)
			}
		}

		return nil
	})

	if err != nil {
		return ret, Error("WalkDir", err)
	}
	return ret, nil
}

func MaybeCreate(folder string) error {
	return os.MkdirAll(folder, 0777)
}

func MoveFile(src, dst string) *Result {
	if src == "" {
		return MsgError("CheckSourcePath", "no source path to move")
	}

	if fileInfo, err := os.Stat(dst); err == nil {
		if fileInfo.IsDir() {
			_, fn := filepath.Split(src)
			dst = filepath.Join(dst, fn)
		}
	}

	dstFolder, _ := filepath.Split(dst)
	if err := MaybeCreate(dstFolder); err != nil {
		return Error("MaybeCreate: "+dstFolder, err)
	}

	inputFile, err := os.Open(src)
	if err != nil {
		return Error("Open: "+src, err)
	}
	outputFile, err := os.Create(dst)
	if err != nil {
		inputFile.Close()
		return Error("Create: "+dst, err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return Error("Copy "+src+" to "+dst, err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(src)
	if err != nil {
		return Error("Remove "+src, err)
	}

	return nil
}
