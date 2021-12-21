package main

import (
	"github.com/SGPractice/glow-bot/client"
	"github.com/SGPractice/glow-bot/command"
	"github.com/bwmarrin/discordgo"
)

func init() {
	command.Register(command.NewCommand("/embed", embed{}))
}

type embed struct{}

func (embed) Run(client *client.Client, msg *discordgo.Message, args []string) {
	// TODO
}
