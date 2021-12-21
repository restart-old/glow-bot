package client

import (
	"github.com/SGPractice/link"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

type Client struct {
	*discordgo.Session
	*link.Linker
	*Config
}

func New(c *Config, l *link.Linker) *Client {
	return &Client{
		Config: c,
		Linker: l,
	}
}

func (client *Client) Start() (err error) {
	client.Session, err = discordgo.New(client.BotToken)
	if err != nil {
		return err
	}
	client.Identify.Intents = discordgo.IntentsAll
	err = client.Session.Open()
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) CloseOnProgramEnd() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	_ = client.Close()
}
