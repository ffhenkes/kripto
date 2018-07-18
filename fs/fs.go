package fs

import (
	"fmt"
	"io/ioutil"
	"os"
)

type (
	FileSystem struct {
		path string
	}
)

func NewFileSystem(path string) *FileSystem {
	return &FileSystem{path}
}

func (fs *FileSystem) Touch(filename string, data []byte) error {

	f, err := os.Create(fmt.Sprintf("%s/%s.secret", fs.path, filename))
	if err != nil {
		return err
	}

	defer f.Close()
	f.Write(data)

	return nil
}

func (fs *FileSystem) Read(filename string) ([]byte, error) {

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.secret", fs.path, filename))
	if err != nil {
		return nil, err
	}

	return data, nil
}
