package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NeowayLabs/logger"
	"github.com/ffhenkes/kripto/algo"
	"github.com/ffhenkes/kripto/fs"
	"github.com/ffhenkes/kripto/model"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"

	"github.com/julienschmidt/httprouter"
)

var logR = logger.Namespace("kripto.router")

const (
	path           = "/data/secrets"
	tmp_passphrase = "avocado"
)

type (
	Router struct {
	}

	Health struct {
		Msg string `json:"msg"`
	}
)

func NewRouter() *Router {
	return &Router{}
}

func (router *Router) Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	health := Health{
		Msg: "I'm alive!",
	}

	h, err := json.Marshal(health)
	if err != nil {
		log.Error("Json parser return with errors: %v", err)
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", h)
}

func (router *Router) CreateSecret(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	sec_request := model.Secret{}

	err := json.NewDecoder(r.Body).Decode(&sec_request)
	if err != nil {
		logR.Error("Decode error: %v", err)
	}

	jsec, err := json.Marshal(sec_request)
	if err != nil {
		logR.Error("Marshal error: %v", err)
	}

	symmetrical := algo.NewSymmetrical()
	cypher, err := symmetrical.Encrypt(jsec, tmp_passphrase)
	if err != nil {
		logR.Error("Encrypt error: %v", err)
	}

	sys := fs.NewFileSystem(path)
	err = sys.Touch(sec_request.App, cypher)
	if err != nil {
		logR.Error("Touch error: %v", err)
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (router *Router) GetSecretsByApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	app := r.URL.Query().Get("app")

	sys := fs.NewFileSystem(path)

	data, err := sys.Read(app)
	if err != nil {
		logR.Error("Read error: %v", err)
	}

	var b []byte
	if len(data) > 0 {

		symmetrical := algo.NewSymmetrical()

		b, err = symmetrical.Decrypt(data, tmp_passphrase)
		if err != nil {
			logR.Error("Decrypt error: %v", err)
		}
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(b))
}
