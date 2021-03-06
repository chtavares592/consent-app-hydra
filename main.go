package main

import (
	"log"

	"github.com/chtavares592/consent-app-hydra/handler"
	"github.com/labstack/echo"
	"github.com/ory/hydra/sdk/go/hydra"
)

func setupTestingHydra() (*hydra.CodeGenSDK, error) {
	client, err := hydra.NewSDK(&hydra.Configuration{
		ClientID:     "auth-server",
		ClientSecret: "auth-secret",
		PublicURL:    "http://localhost:9000",
		AdminURL:     "http://localhost:9001",
		Scopes:       []string{"hydra.clients"},
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

	e.GET("/auth/consent", worker.HandlerConsent)
	e.GET("/auth/login", worker.HandlerLogin)

	e.Logger.Fatal(e.Start(":3000"))
}
