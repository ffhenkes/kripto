package auth

import (
	"fmt"
	"strings"

	"github.com/ffhenkes/kripto/algo"
	"github.com/ffhenkes/kripto/fs"
	"github.com/ffhenkes/kripto/model"
)

const (
	data_authdb = "/data/authdb"
)

type (
	Login struct {
		Credentials *model.Credentials
	}
)

func NewLogin(c *model.Credentials) *Login {
	return &Login{c}
}

func (l *Login) CheckCredentials(phrase string) (bool, error) {

	sys := fs.NewFileSystem(data_authdb)

	data, err := sys.ReadAuth(l.Credentials.Username)
	if err != nil {
		return false, err
	}

	var b []byte
	if len(data) == 0 {
		return false, nil
	}

	symmetrical := algo.NewSymmetrical()

	b, err = symmetrical.Decrypt(data, phrase)
	if err != nil {
		return false, err
	}

	output := strings.Split(string(b), "@")
	username := output[0]
	passwd := output[1]
	hashed_passwd := l.HashPassword()

	if username == l.Credentials.Username && passwd == hashed_passwd {
		return true, nil
	}

	return false, nil
}

func (l *Login) HashPassword() string {
	return fmt.Sprintf("%x", algo.MakeSimpleHash(l.Credentials.Password))
}
