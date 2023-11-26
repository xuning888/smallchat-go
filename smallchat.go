package main

import (
	serverv1 "smallchat/tcp/server/v1"
	serverv2 "smallchat/tcp/server/v2"
)

const maxClients = 2
const listenPort = 7712

func main() {
	testV1()
}

func testV2() {
	server := serverv2.CreateServer("127.0.0.1", listenPort, maxClients)
	server.Spin()
}

func testV1() {
	chatServer := serverv1.CreateServer(
		"127.0.0.1",
		listenPort,
		maxClients,
	)
	chatServer.Spin()
}
