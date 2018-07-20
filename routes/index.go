package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NeowayLabs/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/ffhenkes/kripto/algo"
	"github.com/ffhenkes/kripto/fs"
	"github.com/ffhenkes/kripto/model"

	"github.com/julienschmidt/httprouter"
)

var logR = logger.Namespace("kripto.router")

const (
	data_secrets     = "/data/secrets"
	data_authdb      = "/data/authdb"
	data_rsa         = "/data/rsa"
	private_key_name = "kripto.rsa"
	public_key_name  = "kripto.rsa.pub"
	sign_method      = "RS256"
	time_frame       = 36000
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

	auth := model.Auth{}

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		logR.Error("Decode error: %v", err)
	}

	sys := fs.NewFileSystem(data_authdb)

	data, err := sys.ReadAuth(auth.Username)
	if err != nil {
		logR.Error("Read error: %v", err)
	}

	var b []byte
	if len(data) == 0 {
		msg := map[string]string{"msg": "Invalid username!"}
		m, _ := json.Marshal(msg)

		responseHeader(w, http.StatusBadRequest)
		fmt.Fprintf(w, "%s", m)
		return
	}

	symmetrical := algo.NewSymmetrical()

	b, err = symmetrical.Decrypt(data, router.phrase)
	if err != nil {
		logR.Error("Decrypt error: %v", err)
		responseHeader(w, http.StatusPreconditionFailed)
		return
	}

	output := strings.Split(string(b), "@")
	username := output[0]
	passwd := output[1]
	hashed_passwd := algo.MakeSimpleHash(auth.Password)

	if username == auth.Username && passwd == string(hashed_passwd) {

		rsys := fs.NewFileSystem(data_rsa)

		privateKey, err := rsys.ReadKey(private_key_name)
		if err != nil {
			logR.Error("Read key error: %v", err)
		}

		token := jwt.New(jwt.GetSigningMethod(sign_method))
		token.Claims = &model.CustomClaims{
			&jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
			auth.Username,
		}

		signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
		if err != nil {
			logR.Error("Error parsing key: %v", err)
		}

		tokenString, err := token.SignedString(signKey)
		if err != nil {
			logR.Error("Error signing token: %v", err)
			responseHeader(w, http.StatusExpectationFailed)
			return
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

	authorized, err := checkToken(r)
	if !authorized && err != nil {
		responseHeader(w, http.StatusUnauthorized)
		return
	}

	sec_request := model.Secret{}

	err = json.NewDecoder(r.Body).Decode(&sec_request)
	if err != nil {
		logR.Error("Decode error: %v", err)
	}

	jsec, err := json.Marshal(sec_request)
	if err != nil {
		logR.Error("Marshal error: %v", err)
	}

	symmetrical := algo.NewSymmetrical()
	cypher, err := symmetrical.Encrypt(jsec, router.phrase)
	if err != nil {
		logR.Error("Encrypt error: %v", err)
	}

	sys := fs.NewFileSystem(data_secrets)
	err = sys.MakeSecret(sec_request.App, cypher)
	if err != nil {
		logR.Error("Touch error: %v", err)
	}

	responseHeader(w, http.StatusCreated)
}

// GetSecretsByApp decrypts and returns the required secrets by app
func (router *Router) GetSecretsByApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authorized, err := checkToken(r)
	if err != nil {
		logR.Error("Can't validate token: %v", err)
		responseHeader(w, http.StatusExpectationFailed)
		return
	}

	if !authorized {
		responseHeader(w, http.StatusUnauthorized)
		return
	}

	app := r.URL.Query().Get("app")

	sys := fs.NewFileSystem(data_secrets)

	data, err := sys.ReadSecret(app)
	if err != nil {
		logR.Error("Read error: %v", err)
	}

	var b []byte
	if len(data) > 0 {

		symmetrical := algo.NewSymmetrical()

		b, err = symmetrical.Decrypt(data, router.phrase)
		if err != nil {
			logR.Error("Decrypt error: %v", err)
		}
	}

	responseHeader(w, http.StatusOK)
	fmt.Fprintf(w, "%s", string(b))
}

// RemoveSecretsByApp removes the required secret from the file system by app
func (router *Router) RemoveSecretsByApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	authorized, err := checkToken(r)
	if !authorized && err != nil {
		responseHeader(w, http.StatusUnauthorized)
		return
	}

	app := r.URL.Query().Get("app")

	sys := fs.NewFileSystem(data_secrets)

	err = sys.DeleteSecret(app)
	if err != nil {
		logR.Error("Delete error: %v", err)
	}

	responseHeader(w, http.StatusNoContent)
}

// checkToken utilitary for token validation
func checkToken(r *http.Request) (bool, error) {

	sys := fs.NewFileSystem(data_rsa)

	pub, err := sys.ReadKey(public_key_name)
	if err != nil {
		return false, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return false, err
	}

	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	expires_at := time.Unix(token.Claims.(*model.CustomClaims).StandardClaims.ExpiresAt, 0)
	username := token.Claims.(*model.CustomClaims).Username

	if err != nil {
		logR.Warn("User: %s Non valid token! Expires At: %v", username, expires_at)
		return token.Valid, err
	}

	logR.Info("User: %s Token Expires At: %v", username, expires_at)
	return token.Valid, nil
}

// responseHeader utilitary function to set the output response headers
func responseHeader(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
