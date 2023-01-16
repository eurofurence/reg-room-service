package main

import (
	"fmt"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/logging/consolelogging/logformat"
	"github.com/eurofurence/reg-room-service/web"
	"log"
)

func main() {
	err := config.LoadConfiguration("config.yaml")
	if err != nil {
		log.Fatal(logformat.Logformat("ERROR", "00000000", fmt.Sprintf("Error while loading configuration: %v", err)))
	}
	server := web.Create()
	web.Serve(server)
}
