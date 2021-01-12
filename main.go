package main

import "github.com/eurofurence/reg-room-service/web"

func main() {
	server := web.Create()
	web.Serve(server)
}
