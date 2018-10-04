package main

import (
	"log"

	"github.com/chtavares592/consent-app-hydra/handler"
	"github.com/labstack/echo"
	"github.com/ory/hydra/sdk/go/hydra"
)

func setupTestingHydra() (*hydra.CodeGenSDK, error) {
	// client-ID hydra
	client, err := hydra.NewSDK(&hydra.Configuration{
		ClientID:     "admin",
		ClientSecret: "demo-password",
		EndpointURL:  "http://localhost:9000",
		Scopes:       []string{"hydra.*"},
	})

	return client, err
}

func main() {
	worker := &handler.Worker{}
	var err error
	worker.Client, err = setupTestingHydra()
	if err != nil {
		log.Fatal("Error init hydra sdk")
	}

	e := echo.New()

	e.GET("/consent", worker.HandlerConsent)
	e.GET("/login", worker.HandlerLogin)

	e.Logger.Fatal(e.Start(":3000"))
}