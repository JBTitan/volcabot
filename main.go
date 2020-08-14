package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// discord is the Session being used by the bot.
var discord *discordgo.Session = initSession()

func initSession() *discordgo.Session {
	session, err := discordgo.New()
	if err != nil {
		panic(err)
	}
	session.StateEnabled = true
	return session
}

func main() {
	{
		env, ok := os.LookupEnv("VOLCABOT_LOG_LEVEL")
		if ok {
			lvl, err := log.ParseLevel(env)
			if err != nil {
				panic(err)
			}
			log.SetLevel(lvl)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	}

	{
		token, ok := os.LookupEnv("VOLCABOT_TOKEN")
		if !ok {
			fmt.Fprintln(os.Stderr, "Please pass a bot token using the environment variable VOLCABOT_TOKEN.")
			os.Exit(1)
		}
		discord.Identify.Token = "Bot " + token
		discord.Token = "Bot " + token
	}

	discord.AddHandler(func(s *discordgo.Session, ev *discordgo.Connect) {
		log.Debug("Connected to websocket endpoint")
	})

	err := discord.Open()
	if err != nil {
		panic(err)
	}
	defer discord.Close()
	defer log.Info("Discord bot closing")

	log.Println("Discord bot opened")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-ch
}
