package main

import (
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/web"
	"log"
)

func main() {
	err := config.LoadConfiguration("config.yaml")
	if err != nil {
		log.Fatalf("Error while loading configuration: %v", err)
	}
	server := web.Create()
	web.Serve(server)
}
