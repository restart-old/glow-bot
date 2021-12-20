package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	session, err := discordgo.New("Bot OTIwNDU0NTk4OTcwNDc0NDk3.YbkmJQ.jEsyHsuks7-SnD0gnnLAplnLhXw")
	if err != nil {
		log.Fatalln(err)
	}
	session.Identify.Intents = discordgo.IntentsGuildMembers
	if err := session.Open(); err != nil {
		log.Fatalln(err)
	}
	ranks := map[string]string{
		"Zeus":   "922234208787759125",
		"Kratos": "922234180870492301",
		"Triton": "922234161253716029",
		"Hermes": "922234150130430013",
		"Ares":   "922234108330008656",
	}
	go func() {
		var after string
		ticker := time.NewTicker(1 * time.Minute / 2)
		for {
			<-ticker.C
			members, _ := session.GuildMembers("914172147520401448", after, 1000)
			if len(members) <= 0 {
				after = ""
			}
			for _, m := range members {
				if r, ok, _ := linker.LinkedFromDiscordID(m.User.ID); ok {
					fmt.Println(r.Username())
					if rows, err := linker.DB().Query(fmt.Sprintf("SELECT role FROM playerdata WHERE username='%s';", r.Username())); err == nil {
						var role string
						for rows.Next() {
							rows.Scan(&role)
						}
						rank, ok := ranks[role]
						fmt.Println(role, rank)
						if ok {
							session.GuildMemberRoleAdd("914172147520401448", m.User.ID, rank)
						}
					}
				}
				after = m.User.ID
			}
		}
	}()
	session.AddHandler(MessageCreate)
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	session.Close()
}
