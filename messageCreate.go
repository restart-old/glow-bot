package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SGPractice/link"
	"github.com/bwmarrin/discordgo"
	"github.com/go-sql-driver/mysql"
)

var storer = link.NewJSONStorer("/home/debian/link/")
var linker *link.Linker

func init() {
	config := mysql.NewConfig()
	config.DBName = "GlowHCF"
	config.User = "root"
	config.Addr = ":3306"
	config.Passwd = "f37JZEUm2QFexguhRuyscW{AdrKr86KajFGf%VT2h6BJUUF"
	config.Net = "tcp"

	connector, _ := mysql.NewConnector(config)
	db := sql.OpenDB(connector)
	linker = link.NewLinker(db, storer)
}

func MessageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
	args := strings.Split(msg.Content, " ")
	switch args[0] {
	case "/link":
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if len(args) < 2 {
			go func() {
				m, _ := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%v, Not enought arguments were given.", msg.Author.Mention()))
				time.Sleep(5 * time.Second)
				s.ChannelMessageDelete(m.ChannelID, m.ID)
			}()
			return
		}
		if r, ok, err := linker.LinkedFromDiscordID(msg.Author.ID); !ok {
			if err != nil {
				s.ChannelMessageSend(msg.ChannelID, "An error occurred, please notify staff!")
				return
			}
			l(msg.Author.ID, args[1], msg.Message, s)
		} else {
			go func() {
				m, _ := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%v, You are already linked with the username **%s**, do /unlink in the lobby if you wish to unlink.", msg.Author.Mention(), r.Username()))
				time.Sleep(5 * time.Second)
				s.ChannelMessageDelete(m.ChannelID, m.ID)
			}()
		}
	case "/info":
		if r, ok, err := linker.LinkedFromDiscordID(msg.Author.ID); !ok {
			if err != nil {
				s.ChannelMessageSend(msg.ChannelID, "An error occurred, please notify staff!")
				return
			}
			go func() {
				m, _ := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%v, You are not linked.", msg.Author.Mention()))
				time.Sleep(20 * time.Second)
				s.ChannelMessageDelete(m.ChannelID, m.ID)
			}()
		} else {
			embed := &discordgo.MessageEmbed{}
			embed.Title = fmt.Sprintf("Informations of **%s** (%s)", msg.Author.Username, r.DiscordID())
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Username:",
				Value:  r.Username(),
				Inline: true,
			})
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Linked since:",
				Value:  r.LinkedSince().Format("2006-01-02"),
				Inline: true,
			})
			s.ChannelMessageSendEmbed(msg.ChannelID, embed)
		}
	}
}

func l(discordID, code string, msg *discordgo.Message, s *discordgo.Session) {
	if username, _, ok := linker.LoadByCode(code); ok {
		if err := linker.Link(username, code, discordID); err != nil {
			fmt.Println(err)
			return
		}
		go func() {
			m, _ := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%v, You are now linked with the username **%s**!", msg.Author.Mention(), username))
			time.Sleep(5 * time.Second)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}()
	} else {
		go func() {
			m, _ := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%v, Your code is invalid!", msg.Author.Mention()))
			time.Sleep(5 * time.Second)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}()
	}
}
