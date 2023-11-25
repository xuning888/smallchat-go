package main

import (
	"smallchat/tcp/server"
)

const maxClients = 100
const listenPort = 7712

func main() {
	server := server.CreateServer(
		"127.0.0.1",
		listenPort,
		maxClients,
	)
	server.Spin()
}
