package auth

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ffhenkes/kripto/algo"
	"github.com/ffhenkes/kripto/fs"
	"github.com/ffhenkes/kripto/model"
)

const (
	dataAuthdb = "/data/authdb"
)

type (
	// Login is used to create and validate Credentials
	Login struct {
		Credentials *model.Credentials
	}
)

// NewLogin returns a Login type with embed Credentials
func NewLogin(c *model.Credentials) *Login {
	return &Login{c}
}

// AddCredentials creates a new user file on the disk containing user and password data encrypted using the kripto built in passphrase
func (l *Login) AddCredentials(phrase string) error {

	passwd := l.HashPassword()
	userString := fmt.Sprintf("%s@%s@%d", l.Credentials.Username, passwd, l.Credentials.TokenExpiresIn)

	symmetrical := algo.NewSymmetrical()
	data, err := symmetrical.Encrypt([]byte(userString), phrase)
	if err != nil {
		return err
	}

	sys := fs.NewFileSystem(dataAuthdb)
	err = sys.MakeAuth(l.Credentials.Username, data)
	return err
}

// CheckCredentials retrieve the user data from file, decrypt it and returns a boolean sign
func (l *Login) CheckCredentials(phrase string) (bool, error) {

	sys := fs.NewFileSystem(dataAuthdb)

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
	hashedPasswd := l.HashPassword()

	t, err := strconv.Atoi(output[2])
	if err != nil {
		return false, err
	}

	l.Credentials.TokenExpiresIn = time.Duration(t)

	if username == l.Credentials.Username && passwd == hashedPasswd {
		return true, nil
	}

	return false, nil
}

// HashPassword create a string hash using sha256 built within the MakeSimpleHash algorithm
func (l *Login) HashPassword() string {
	return fmt.Sprintf("%x", algo.MakeSimpleHash(l.Credentials.Password))
}
