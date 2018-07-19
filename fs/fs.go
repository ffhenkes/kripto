package fs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/NeowayLabs/logger"
)

type (
	FileSystem struct {
		path string
	}
)

func NewFileSystem(path string) *FileSystem {
	return &FileSystem{path}
}

func (fs *FileSystem) MakeAuth(filename string, data []byte) error {
	logger.Info("%s", filename)
	err := touch(authdb(fs.path, filename), data)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileSystem) ReadAuth(filename string) ([]byte, error) {

	data, err := read(authdb(fs.path, filename))
	if err != nil {
		return nil, err
	}

	return data, err
}

func (fs *FileSystem) DeleteAuth(filename string) error {

	err := del(authdb(fs.path, filename))
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileSystem) MakeSecret(filename string, data []byte) error {
	err := touch(secret(fs.path, filename), data)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileSystem) ReadSecret(filename string) ([]byte, error) {

	data, err := read(secret(fs.path, filename))
	if err != nil {
		return nil, err
	}

	return data, err
}

func (fs *FileSystem) DeleteSecret(filename string) error {

	err := del(secret(fs.path, filename))
	if err != nil {
		return err
	}

	return nil
}

// helpers
func touch(out string, data []byte) error {

	f, err := os.Create(out)
	if err != nil {
		return err
	}

	defer f.Close()
	f.Write(data)

	return nil
}

func read(out string) ([]byte, error) {

	data, err := ioutil.ReadFile(out)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func del(out string) error {

	err := os.Remove(out)
	if err != nil {
		return err
	}

	return nil
}

func authdb(p, f string) string {
	return fmt.Sprintf("%s/%s.auth", p, f)
}

func secret(p, f string) string {
	return fmt.Sprintf("%s/%s.secret", p, f)
}
