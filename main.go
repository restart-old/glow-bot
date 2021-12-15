package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	session, err := discordgo.New("Bot OTIwNDU0NTk4OTcwNDc0NDk3.YbkmJQ.jEsyHsuks7-SnD0gnnLAplnLhXw")
	if err != nil {
		log.Fatalln(err)
	}
	if err := session.Open(); err != nil {
		log.Fatalln(err)
	}
	session.AddHandler(MessageCreate)
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	session.Close()
}
