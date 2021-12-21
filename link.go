package main

import (
	"fmt"
	"github.com/SGPractice/glow-bot/client"
	"github.com/SGPractice/glow-bot/command"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func init() {
	command.Register(command.NewCommand("/link", linkCMD{}))
}

type linkCMD struct{}

func (linkCMD) Run(client *client.Client, msg *discordgo.Message, args []string) {
	// Delete the msg
	if err := client.ChannelMessageDelete(msg.ChannelID, msg.ID); err != nil {
		log.Println(err)
		return
	}

	// Check if there are enough arguments
	if len(args) < 2 {
		command.SendTempMsg(client, msg.ChannelID, fmt.Sprintf("%s, Not enough arguments provided.", msg.Author.Mention()), 5*time.Second)
		return
	}

	// Check if the author is already linked
	if r, ok, err := client.LinkedFromDiscordID(msg.Author.ID); !ok {
		if err != nil {
			command.SendTempMsg(client, msg.ChannelID, fmt.Sprintf("%s, An error occured, please notify staff.", msg.Author.Mention()), 10*time.Second)
			return
		}
		// Check if the code is valid
		if username, _, ok := client.LoadByCode(args[1]); ok {
			if err = client.Link(username, args[1], msg.Author.ID); err != nil {
				log.Println(err)
				return
			}
			command.SendTempMsg(client, msg.ChannelID, fmt.Sprintf("%v, You are now linked with the username **%s**!", msg.Author.Mention(), username), 5*time.Second)
			// Not valid
		} else {
			command.SendTempMsg(client, msg.ChannelID, fmt.Sprintf("%v, Your code is invalid!", msg.Author.Mention()), 5*time.Second)
		}
	} else {
		command.SendTempMsg(client, msg.ChannelID, fmt.Sprintf("%v, You are already linked with the username **%s**, do /unlink in the lobby if you wish to unlink.", msg.Author.Mention(), r.Username()), 5*time.Second)
	}
}
