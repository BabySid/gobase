package gobase

import (
	"io/fs"
	"io/ioutil"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriteFile(filename string, data []byte, perm fs.FileMode) error {
	tmp := filename + ".tmp"
	err := ioutil.WriteFile(tmp, data, perm)
	if err != nil {
		return err
	}

	return os.Rename(tmp, filename)
}
