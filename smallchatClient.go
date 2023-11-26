package main

import (
	"smallchat/tcp/client"
	"smallchat/tcp/client/cmd"
)

func main() {
	chatClient := client.CreateChatClient("127.0.0.1", 7712)
	cmds := commands()
	for cmd, command := range cmds {
		chatClient.RegisterCmd(cmd, command)
	}
	chatClient.Spin()
}

func commands() map[string]cmd.Command {
	cmds := make(map[string]cmd.Command)
	cmds["/help"] = cmd.CreateHelp(cmds)
	cmds["/exit"] = &cmd.Exit{}
	return cmds
}
