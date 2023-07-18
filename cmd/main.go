package main

import (
	"github.com/eurofurence/reg-room-service/internal/web/app"
	"os"
)

func main() {
	os.Exit(app.New().Run())
}
