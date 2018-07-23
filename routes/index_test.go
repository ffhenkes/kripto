package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ffhenkes/kripto/auth"
	"github.com/ffhenkes/kripto/fs"
	"github.com/ffhenkes/kripto/model"
)

const (
	test_passphrase  = "avocado"
	test_data_authdb = "/data/authdb"
	test_user        = "ffhenkes"
	test_passwd      = "test"
	bad_password     = "penguim"
	bad_username     = "jonah"
)

var c *model.Credentials
var s *model.Secret
var token string = ""

func TestHealth(t *testing.T) {

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/health", nil)

	router := NewRouter(test_passphrase)
	router.Health(res, req, nil)

	status := res.Code

	if status != http.StatusOK {
		t.Errorf("Bad status, unhealthy! Got %v expected %v", status, http.StatusOK)
	}
}

func TestShouldLogin(t *testing.T) {

	err := before()
	if err != nil {
		t.Fatal(err)
	}

	jc, _ := json.Marshal(c)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/authenticate", bytes.NewReader(jc))

	router := NewRouter(test_passphrase)
	router.Authenticate(res, req, nil)

	status := res.Code

	if status != http.StatusCreated {
		t.Errorf("Bad status! Got %v expected %v", status, http.StatusCreated)
	}

	err = tearDown()
	if err != nil {
		t.Fatal(err)
	}
}

func TestShouldNotLogin(t *testing.T) {

	err := before()
	if err != nil {
		t.Fatal(err)
	}

	c.Password = bad_password

	jc, _ := json.Marshal(c)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/authenticate", bytes.NewReader(jc))

	router := NewRouter(test_passphrase)
	router.Authenticate(res, req, nil)

	status := res.Code

	if status != http.StatusUnauthorized {
		t.Errorf("Bad status! Got %v expected %v", status, http.StatusUnauthorized)
	}

	err = tearDown()
	if err != nil {
		t.Fatal(err)
	}
}

func TestShouldNotFindUser(t *testing.T) {

	err := before()
	if err != nil {
		t.Fatal(err)
	}

	c.Username = bad_username

	jc, _ := json.Marshal(c)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/authenticate", bytes.NewReader(jc))

	router := NewRouter(test_passphrase)
	router.Authenticate(res, req, nil)

	status := res.Code

	if status != http.StatusInternalServerError {
		t.Errorf("Bad status! Got %v expected %v", status, http.StatusInternalServerError)
	}

	err = tearDown()
	if err != nil {
		t.Fatal(err)
	}

}

func TestShouldCreateSecret(t *testing.T) {

	err := before()
	if err != nil {
		t.Fatal(err)
	}

	j := auth.NewJwtAuth(c)

	token, err = j.GenerateToken()
	if err != nil {
		t.Fatal(err)
	}

	s = &model.Secret{
		App: "kripto_test",
		Vars: map[string]string{
			"some_url":    "http://something.krp",
			"some_passwd": "xyzptbo2294%@",
		},
	}

	jsec, _ := json.Marshal(s)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/secrets", bytes.NewReader(jsec))
	req.Header.Add("Authorization", token)

	router := NewRouter(test_passphrase)

	router.CreateSecret(res, req, nil)

	status := res.Code

	if status != http.StatusCreated {
		t.Errorf("Bad status! Got %v expected %v", status, http.StatusCreated)
	}

}

func TestShouldGetSecret(t *testing.T) {

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/secrets?app=kripto_test", nil)
	req.Header.Add("Authorization", token)

	router := NewRouter(test_passphrase)

	router.GetSecretsByApp(res, req, nil)

	status := res.Code

	cracked, err := decodeSecret(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if status != http.StatusOK {
		t.Errorf("Bad status! Got %v expected %v", status, http.StatusOK)
	}

	match := reflect.DeepEqual(s, cracked)
	if !match {
		t.Errorf("Secret not cracked! %t", match)
	}
}

func TestShouldDelSecret(t *testing.T) {

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/v1/secrets?app=kripto_test", nil)
	req.Header.Add("Authorization", token)

	router := NewRouter(test_passphrase)

	router.RemoveSecretsByApp(res, req, nil)

	status := res.Code

	if status != http.StatusNoContent {
		t.Errorf("Bad status! Got %v expected %v", status, http.StatusNoContent)
	}

	err := tearDown()
	if err != nil {
		t.Fatal(err)
	}

}

func before() error {

	c = &model.Credentials{
		Username: test_user,
		Password: test_passwd,
	}

	l := auth.NewLogin(c)
	err := l.AddCredentials(test_passphrase)
	if err != nil {
		return err
	}

	return nil
}

func tearDown() error {

	sys := fs.NewFileSystem(test_data_authdb)

	err := sys.RemovePath()
	if err != nil {
		return err
	}

	return nil
}

func decodeSecret(r io.Reader) (*model.Secret, error) {
	var secret *model.Secret
	if err := json.NewDecoder(r).Decode(&secret); err != nil {
		return nil, errors.New("Bad decode: " + err.Error())
	}
	return secret, nil
}
