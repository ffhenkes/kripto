package fs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type (
	// FileSystem represent a type that loads operations that can be performed into the file system.
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

	s, err := sanitize(filename)
	if err != nil {
		return err
	}

	err = mkdir(fs.path)
	if err != nil {
		return err
	}
	err = touch(authdb(fs.path, s), data)
	return err
}

// ReadAuth reads the authdb refered file
func (fs *FileSystem) ReadAuth(filename string) ([]byte, error) {

	s, err := sanitize(filename)
	if err != nil {
		return nil, err
	}

	data, err := read(authdb(fs.path, s))
	if err != nil {
		return nil, err
	}

	return data, err
}

// DeleteAuth removes the specific auth file
func (fs *FileSystem) DeleteAuth(filename string) error {

	s, err := sanitize(filename)
	if err != nil {
		return err
	}

	err = del(authdb(fs.path, s))
	return err
}

// ReadKey reads the rsa private key from rsa directory
func (fs *FileSystem) ReadKey(keyname string) ([]byte, error) {

	s, err := sanitize(keyname)
	if err != nil {
		return nil, err
	}

	key, err := read(rsa(fs.path, s))
	if err != nil {
		return nil, err
	}

	return key, nil
}

// ReadPublicKey reads the rsa public key from rsa directory
func (fs *FileSystem) ReadPublicKey(keyname string) ([]byte, error) {

	s, err := sanitize(keyname)
	if err != nil {
		return nil, err
	}

	key, err := read(rsapub(fs.path, s))
	if err != nil {
		return nil, err
	}

	return key, nil
}

// MakeSecret creates a new secret file into the secrets directory
func (fs *FileSystem) MakeSecret(filename string, data []byte) error {

	s, err := sanitize(filename)
	if err != nil {
		return err
	}

	err = mkdir(fs.path)
	if err != nil {
		return err
	}

	err = touch(secret(fs.path, s), data)
	return err
}

// ReadSecret creates reads a specific secret from the secrets directory
func (fs *FileSystem) ReadSecret(filename string) ([]byte, error) {

	s, err := sanitize(filename)
	if err != nil {
		return nil, err
	}

	data, err := read(secret(fs.path, s))
	if err != nil {
		return nil, err
	}

	return data, err
}

// DeleteSecret removes a specific secret file
func (fs *FileSystem) DeleteSecret(filename string) error {

	s, err := sanitize(filename)
	if err != nil {
		return err
	}

	err = del(secret(fs.path, s))
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

	defer closeFile(f)

	_, err = f.Write(data)
	return err
}

func read(out string) ([]byte, error) {

	// the annotation below suppress gosec warning
	// this particular case is solved by the sanitize function
	/* #nosec */
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

func sanitize(input string) (string, error) {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}

	s := reg.ReplaceAllString(input, "")
	if "" == s {
		return "", errors.New("sanitize: empty string")
	}

	return s, nil
}

func rsa(p, f string) string {
	return fmt.Sprintf("%s/%s.rsa", p, f)
}

func rsapub(p, f string) string {
	return fmt.Sprintf("%s/%s.rsa.pub", p, f)
}

func authdb(p, f string) string {
	return fmt.Sprintf("%s/.%s.auth", p, f)
}

func secret(p, f string) string {
	return fmt.Sprintf("%s/%s.secret", p, f)
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		panic(err)
	}
}
