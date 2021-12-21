package command

import (
	"github.com/SGPractice/glow-bot/client"
	"github.com/bwmarrin/discordgo"
	"strings"
	"sync"
	"time"
)

var cMutex sync.Mutex
var commands []*Command

func Register(c *Command) {
	cMutex.Lock()
	defer cMutex.Unlock()
	commands = append(commands, c)
}

func NewCommand(name string, runnable Runnable) *Command {
	return &Command{Runnable: runnable, name: name}
}

type Command struct {
	name string
	Runnable
}

type Runnable interface {
	Run(client *client.Client, msg *discordgo.Message, args []string)
}

func StartHandlingCommands(c *client.Client) {
	f := func(s *discordgo.Session, msg *discordgo.MessageCreate) {
		args := strings.Split(msg.Content, " ")
		for _, cmd := range commands {
			if cmd.name == args[0] {
				cmd.Run(c, msg.Message, args)
			}
		}
	}
	c.AddHandler(f)
}

func SendTempMsg(client *client.Client, channelID string, msg string, d time.Duration) {
	go func() {
		m, err := client.ChannelMessageSend(channelID, msg)
		if err != nil {
			return
		}
		time.Sleep(d)
		_ = client.ChannelMessageDelete(m.ChannelID, m.ID)
	}()
}
