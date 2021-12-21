package main

import (
	"database/sql"
	"fmt"
	"github.com/RestartFU/gophig"
	client2 "github.com/SGPractice/glow-bot/client"
	"github.com/SGPractice/glow-bot/command"
	"github.com/SGPractice/link"
	mysql2 "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func main() {
	// mysql
	var mysql mysql2.Config
	if err := gophig.GetConf("./mysql", "toml", &mysql); err != nil {
		log.Fatalln(err)
	}

	// connector
	connector, err := mysql2.NewConnector(&mysql)
	if err != nil {
		log.Fatalln(err)
	}
	// linker
	storer := link.NewJSONStorer("/home/debian/link/")
	db := sql.OpenDB(connector)
	linker := link.NewLinker(db, storer)

	// config
	var config client2.Config
	if err = gophig.GetConf("./config", "toml", &config); err != nil {
		log.Fatalln(err)
	}

	// client
	client := client2.New(&config, linker)

	if err = client.Start(); err != nil {
		log.Fatalln(err)
	}
	log.Println("bot started")
	command.StartHandlingCommands(client)

	go func() {
		ranks := map[string]string{
			"Zeus":   "922234208787759125",
			"Kratos": "922234180870492301",
			"Triton": "922234161253716029",
			"Hermes": "922234150130430013",
			"Ares":   "922234108330008656",
		}
		var after string
		ticker := time.NewTicker(1 * time.Minute / 2)
		for {
			<-ticker.C
			members, _ := client.GuildMembers("914172147520401448", after, 1000)
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
							client.GuildMemberRoleAdd("914172147520401448", m.User.ID, rank)
						}
					}
				}
				after = m.User.ID
			}
		}
	}()

	client.CloseOnProgramEnd()
}
