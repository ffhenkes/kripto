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
	testPassphrase = "avocado"
	testDataAuthdb = "/data/authdb"
	testUser       = "ffhenkes"
	testPasswd     = "test"
	badPassword    = "penguim"
	badUsername    = "jonah"
)

var c *model.Credentials
var s *model.Secret
var token string

func TestHealth(t *testing.T) {

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/health", nil)

	router := NewRouter(testPassphrase)
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

	router := NewRouter(testPassphrase)
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

	c.Password = badPassword

	jc, _ := json.Marshal(c)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/authenticate", bytes.NewReader(jc))

	router := NewRouter(testPassphrase)
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

	c.Username = badUsername

	jc, _ := json.Marshal(c)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/authenticate", bytes.NewReader(jc))

	router := NewRouter(testPassphrase)
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

	router := NewRouter(testPassphrase)

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

	router := NewRouter(testPassphrase)

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

	router := NewRouter(testPassphrase)

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
		Username: testUser,
		Password: testPasswd,
	}

	l := auth.NewLogin(c)
	err := l.AddCredentials(testPassphrase)
	return err
}

func tearDown() error {

	sys := fs.NewFileSystem(testDataAuthdb)

	err := sys.RemovePath()
	return err
}

func decodeSecret(r io.Reader) (*model.Secret, error) {
	var secret *model.Secret
	if err := json.NewDecoder(r).Decode(&secret); err != nil {
		return nil, errors.New("Bad decode: " + err.Error())
	}
	return secret, nil
}
