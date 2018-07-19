package main

import (
	"net/http"
	"os"

	"github.com/NeowayLabs/logger"
	"github.com/ffhenkes/kripto/routes"
	"github.com/julienschmidt/httprouter"
)

func main() {

	var logH = logger.Namespace("kripto")

	var addr = os.Getenv("KRIPTO_ADDRESS")
	var phrase = os.Getenv("PHRASE")
	var crt = os.Getenv("CRT_PATH")
	var key = os.Getenv("KEY_PATH")

	// Instantiate a new router
	r := httprouter.New()

	nr := routes.NewRouter(phrase)

	// health check
	r.GET("/v1", nr.Health)
	r.POST("/v1/secrets", nr.CreateSecret)
	r.GET("/v1/secrets", nr.GetSecretsByApp)
	r.DELETE("/v1/secrets", nr.RemoveSecretsByApp)

	logH.Info("Running on %s", addr)

	if err := http.ListenAndServeTLS(addr, crt, key, r); err != nil {
		logH.Fatal("ListenAndServe: %s", err)
	}
}
