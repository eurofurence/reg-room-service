package main

import (
	"os"

	"github.com/eurofurence/reg-room-service/internal/web/app"
)

func main() {
	os.Exit(app.New().Run())
}
