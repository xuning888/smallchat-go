package main

import (
	"smallchat/tcp/server"
)

const maxClients = 2
const listenPort = 7712

func main() {
	chatServer := server.CreateServer(
		"127.0.0.1",
		listenPort,
		maxClients,
	)
	chatServer.Spin()
}
