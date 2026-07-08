package main

import (
	"log"

	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/app"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/server"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("init app failed: %v", err)
	}

	srv := server.New(application)
	if err := srv.Run(); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}
