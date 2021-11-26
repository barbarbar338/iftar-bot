package main

import (
	"fmt"
	"iftarbot/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	language "github.com/moemoe89/go-localization"
)

var (
	token     string
	prefix    string
	guildID   string
	channelID string
	ownerID   string
	data      []Vakit
	isPlaying bool
	i18n 	  *language.Config
)

func main() {
	var err error
	data, err = fetchData()
	if err != nil {
		logger.WithFields(logger.Fields{"component": "main", "action": "fetching data from server."}).
			Errorf("An error occurred while loading dataset, %v", err)

		return
	}

	discord, err := discordgo.New(fmt.Sprintf("Bot %v", token))
	if err != nil {
		logger.WithFields(logger.Fields{"component": "main", "action": "create new discord session instance."}).
			Errorf("An error occurred while creating discord session, %v", err)

		return
	}

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		logger.WithFields(logger.Fields{"component": "main", "action": "create new discord session instance."}).
			Errorf("An error occurred while connecting to discord API, %v", err)

		return
	}

	logger.WithFields(logger.Fields{"component": "main", "action": "create new discord session instance."}).
		Infof("Logged in as %v \n", discord.State.User.Username)

	setupJob(discord)

	sc := make(chan os.Signal, 1)
	signal.Notify(
		sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	<-sc

	logger.WithFields(logger.Fields{"component": "main", "action": "closing discord session instance."}).
		Infof("Closing bot, see you later...")

	err = discord.Close()
	if err != nil {
		logger.WithFields(logger.Fields{"component": "main", "action": "closing discord session instance."}).
			Errorf("An error occurred while closing discord session instance, %v", err)
	}
}
