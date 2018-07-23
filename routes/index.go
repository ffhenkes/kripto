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
	data_secrets = "/data/secrets"
)

type (
	Router struct {
		phrase string
	}

	Health struct {
		Msg string `json:"msg"`
	}
)

func NewRouter(phrase string) *Router {
	return &Router{phrase}
}

// Health is a simple health check to verify the basic app running state
func (router *Router) Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	health := Health{
		Msg: "I'm alive!",
	}

	h, err := json.Marshal(health)
	if err != nil {
		logR.Error("Json parser return with errors: %v", err)
	}

	responseHeader(w, http.StatusOK)
	fmt.Fprintf(w, "%s", h)
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
	}

	if ok {
		jtoken := auth.NewJwtAuth(&c)

		tokenString, err := jtoken.GenerateToken()
		if err != nil {
			serverError(w, err)
		}

		msg := map[string]string{"token": tokenString}
		m, _ := json.Marshal(msg)

		responseHeader(w, http.StatusCreated)
		fmt.Fprintf(w, "%s", m)
		return

	}

	responseHeader(w, http.StatusUnauthorized)
}

// CreateSecret records the requested secrets of an app into file system encripting those with a symmetrical algorithm
func (router *Router) CreateSecret(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authorization := r.Header.Get("Authorization")

	authorized, err := auth.ValidateToken(authorization)
	if !authorized && err != nil {
		responseHeader(w, http.StatusUnauthorized)
		return
	}

	sec_request := model.Secret{}

	err = json.NewDecoder(r.Body).Decode(&sec_request)
	if err != nil {
		serverError(w, err)
	}

	jsec, err := json.Marshal(sec_request)
	if err != nil {
		serverError(w, err)
	}

	symmetrical := algo.NewSymmetrical()
	cypher, err := symmetrical.Encrypt(jsec, router.phrase)
	if err != nil {
		serverError(w, err)
	}

	sys := fs.NewFileSystem(data_secrets)
	err = sys.MakeSecret(sec_request.App, cypher)
	if err != nil {
		serverError(w, err)
	}

	responseHeader(w, http.StatusCreated)
}

// GetSecretsByApp decrypts and returns the required secrets by app
func (router *Router) GetSecretsByApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authorization := r.Header.Get("Authorization")

	authorized, err := auth.ValidateToken(authorization)
	if !authorized && err != nil {
		responseHeader(w, http.StatusUnauthorized)
		return
	}

	app := r.URL.Query().Get("app")

	sys := fs.NewFileSystem(data_secrets)

	data, err := sys.ReadSecret(app)
	if err != nil {
		serverError(w, err)
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
	fmt.Fprintf(w, "%s", string(b))
}

// RemoveSecretsByApp removes the required secret from the file system by app
func (router *Router) RemoveSecretsByApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authorization := r.Header.Get("Authorization")

	authorized, err := auth.ValidateToken(authorization)
	if !authorized && err != nil {
		responseHeader(w, http.StatusUnauthorized)
		return
	}

	app := r.URL.Query().Get("app")

	sys := fs.NewFileSystem(data_secrets)

	err = sys.DeleteSecret(app)
	if err != nil {
		serverError(w, err)
	}

	responseHeader(w, http.StatusNoContent)
}

// serverError utilitary to log the specific server problem and returns 500
func serverError(w http.ResponseWriter, err error) {
	logR.Error("Server error %v", err)
	responseHeader(w, http.StatusInternalServerError)
	return
}

// responseHeader utilitary function to set the output response headers
func responseHeader(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
