package main

import (
	"net/http"
	"os"

	"github.com/NeowayLabs/logger"
	"github.com/ffhenkes/kripto/routes"
	"github.com/julienschmidt/httprouter"
)

var (
	// Phrase is loaded in build time with the encryption key
	Phrase string
)

func main() {

	var (
		logH = logger.Namespace("kripto")
		addr = os.Getenv("KRIPTO_ADDRESS")
		crt  = os.Getenv("CRT_PATH")
		key  = os.Getenv("KEY_PATH")
	)

	// Instantiate a new router
	r := httprouter.New()

	nr := routes.NewRouter(Phrase)

	// health check
	r.GET("/v1/health", nr.Health)
	r.POST("/v1/authenticate", nr.Authenticate)
	r.POST("/v1/secrets", nr.CreateSecret)
	r.GET("/v1/secrets", nr.GetSecretsByApp)
	r.DELETE("/v1/secrets", nr.RemoveSecretsByApp)

	logH.Info("Running on %s", addr)

	if err := http.ListenAndServeTLS(addr, crt, key, r); err != nil {
		logH.Fatal("ListenAndServe: %s", err)
	}
}
