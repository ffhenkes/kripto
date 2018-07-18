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

	// Instantiate a new router
	r := httprouter.New()

	nr := routes.NewRouter()

	// health check
	r.GET("/v1", nr.Health)
	r.POST("/v1/secrets", nr.CreateSecret)
	r.GET("/v1/secrets", nr.GetSecretsByApp)

	logH.Info("Running on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		logH.Fatal("ListenAndServe: %s", err)
	}
}
