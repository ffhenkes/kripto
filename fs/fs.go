package fs

import (
	"fmt"
	"io/ioutil"
	"os"
)

type (
	// FileSystem reprsents the operations that can be performed into the file system.
	// Such as create, read, delete
	FileSystem struct {
		path string
	}
)

// NewFileSystem returns a new FileSystem reference with the embedded base path
func NewFileSystem(path string) *FileSystem {
	return &FileSystem{path}
}

// MakeAuth creates a file into the authdb directory that contains users (credentials)
func (fs *FileSystem) MakeAuth(filename string, data []byte) error {

	err := mkdir(fs.path)
	if err != nil {
		return err
	}

	err = touch(authdb(fs.path, filename), data)
	return err
}

// ReadAuth reads the authdb refered file
func (fs *FileSystem) ReadAuth(filename string) ([]byte, error) {

	data, err := read(authdb(fs.path, filename))
	if err != nil {
		return nil, err
	}

	return data, err
}

// DeleteAuth removes the specific auth file
func (fs *FileSystem) DeleteAuth(filename string) error {

	err := del(authdb(fs.path, filename))
	return err
}

// ReadKey reads the rsa key (public or private) from rsa directory
func (fs *FileSystem) ReadKey(keyname string) ([]byte, error) {

	key, err := read(rsa(fs.path, keyname))
	if err != nil {
		return nil, err
	}

	return key, nil
}

// MakeSecret creates a new secret file into the secrets directory
func (fs *FileSystem) MakeSecret(filename string, data []byte) error {

	err := mkdir(fs.path)
	if err != nil {
		return err
	}

	err = touch(secret(fs.path, filename), data)
	return err
}

// ReadSecret creates reads a specific secret from the secrets directory
func (fs *FileSystem) ReadSecret(filename string) ([]byte, error) {

	data, err := read(secret(fs.path, filename))
	if err != nil {
		return nil, err
	}

	return data, err
}

// DeleteSecret removes a specific secret file
func (fs *FileSystem) DeleteSecret(filename string) error {

	err := del(secret(fs.path, filename))
	return err
}

// RemovePath drops the base path
func (fs *FileSystem) RemovePath() error {

	err := os.RemoveAll(fs.path)
	return err
}

// helpers
func mkdir(path string) error {

	err := os.MkdirAll(path, os.ModePerm)
	return err
}

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
	return err
}

func rsa(p, f string) string {
	return fmt.Sprintf("%s/%s", p, f)
}

func authdb(p, f string) string {
	return fmt.Sprintf("%s/.%s.auth", p, f)
}

func secret(p, f string) string {
	return fmt.Sprintf("%s/%s.secret", p, f)
}
