package gobase

import (
	"bufio"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
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

func ReadLine(filename string, handle func(string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	var data string
	for {
		data, err = r.ReadString('\n')
		data = strings.TrimSpace(data)
		if err = handle(data); err != nil {
			return err
		}

		if err != nil && err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
	}
}
