package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NeowayLabs/logger"
	"github.com/ffhenkes/kripto/algo"
	"github.com/ffhenkes/kripto/auth"
	"github.com/ffhenkes/kripto/fs"
	"github.com/ffhenkes/kripto/model"

	"github.com/julienschmidt/httprouter"
)

var logR = logger.Namespace("kripto.router")

const (
	dataSecrets = "/data/secrets"
)

type (
	// Router represents the http api router that embed the built in passphrase for encryption
	Router struct {
		phrase string
	}
)

// NewRouter returns an http Router reference with the embedded kripto built in passphrase
func NewRouter(phrase string) *Router {
	return &Router{phrase}
}

// Health is a simple health check to verify the basic app running state
func (router *Router) Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	health := map[string]string{
		"msg": "I'm alive!",
	}

	h, err := json.Marshal(health)
	if err != nil {
		logR.Error("Json parser return with errors: %v", err)
	}

	responseHeader(w, http.StatusOK)
	_, err = fmt.Fprintf(w, "%s", h)
	if err != nil {
		logR.Fatal("Bad output: %v", err)
	}
}

// Authenticate is a method for validating user and password returning a signed JWT with 24h expiration time
func (router *Router) Authenticate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	c := model.Credentials{}

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logR.Error("Decode error: %v", err)
	}

	login := auth.NewLogin(&c)

	ok, err := login.CheckCredentials(router.phrase)
	if err != nil {
		serverError(w, err)
		return
	}

	if ok {
		jtoken := auth.NewJwtAuth(&c)

		tokenString, err := jtoken.GenerateToken()
		if err != nil {
			serverError(w, err)
			return
		}

		msg := map[string]string{"token": tokenString}
		m, err := json.Marshal(msg)
		if err != nil {
			serverError(w, err)
			return
		}

		responseHeader(w, http.StatusCreated)
		_, err = fmt.Fprintf(w, "%s", m)
		if err != nil {
			logR.Fatal("Bad output: %v", err)
		}
		return

	}

	responseHeader(w, http.StatusUnauthorized)
}

// CreateSecret records the requested secrets of an app into file system encripting those with a symmetrical algorithm
func (router *Router) CreateSecret(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authorization := r.Header.Get("Authorization")

	authorized, err := auth.ValidateToken(authorization)
	if !authorized && err != nil {
		unauthorized(w, err)
		return
	}

	secRequest := model.Secret{}

	err = json.NewDecoder(r.Body).Decode(&secRequest)
	if err != nil {
		serverError(w, err)
		return
	}

	jsec, err := json.Marshal(secRequest)
	if err != nil {
		serverError(w, err)
		return
	}

	symmetrical := algo.NewSymmetrical()
	cypher, err := symmetrical.Encrypt(jsec, router.phrase)
	if err != nil {
		serverError(w, err)
		return
	}

	sys := fs.NewFileSystem(dataSecrets)
	err = sys.MakeSecret(secRequest.App, cypher)
	if err != nil {
		serverError(w, err)
		return
	}

	responseHeader(w, http.StatusCreated)
}

// GetSecretsByApp decrypts and returns the required secrets by app
func (router *Router) GetSecretsByApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authorization := r.Header.Get("Authorization")

	authorized, err := auth.ValidateToken(authorization)
	if !authorized && err != nil {
		unauthorized(w, err)
		return
	}

	app := r.URL.Query().Get("app")

	sys := fs.NewFileSystem(dataSecrets)

	data, err := sys.ReadSecret(app)
	if err != nil {
		serverError(w, err)
		return
	}

	var b []byte
	if len(data) > 0 {

		symmetrical := algo.NewSymmetrical()

		b, err = symmetrical.Decrypt(data, router.phrase)
		if err != nil {
			serverError(w, err)
		}
	}

	responseHeader(w, http.StatusOK)
	_, err = fmt.Fprintf(w, "%s", string(b))
	if err != nil {
		logR.Fatal("Bad output: %v", err)
	}
}

// RemoveSecretsByApp removes the required secret from the file system by app
func (router *Router) RemoveSecretsByApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authorization := r.Header.Get("Authorization")

	authorized, err := auth.ValidateToken(authorization)
	if !authorized && err != nil {
		unauthorized(w, err)
		return
	}

	app := r.URL.Query().Get("app")

	sys := fs.NewFileSystem(dataSecrets)

	err = sys.DeleteSecret(app)
	if err != nil {
		serverError(w, err)
		return
	}

	responseHeader(w, http.StatusNoContent)
}

// unauthorized utilitary to log the specific auth fail and returns 401
func unauthorized(w http.ResponseWriter, err error) {
	logR.Error("Unauthorized %v", err)
	responseHeader(w, http.StatusUnauthorized)
}

// serverError utilitary to log the specific server problem and returns 500
func serverError(w http.ResponseWriter, err error) {
	logR.Error("Server error %v", err)
	responseHeader(w, http.StatusInternalServerError)
}

// responseHeader utilitary function to set the output response headers
func responseHeader(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
