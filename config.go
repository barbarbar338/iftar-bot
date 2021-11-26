package main

import (
	"fmt"
	"iftarbot/logger"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("An error occured while loading .env file, error : %v \n", err.Error()))
	}

	initLogger() // initialize logger

	token = os.Getenv("BOT_TOKEN")
	prefix = os.Getenv("BOT_PREFIX")
	guildID = os.Getenv("GUILD_ID")
	channelID = os.Getenv("CHANNEL_ID")
	ownerID = os.Getenv("OWNER_ID")
}


func initLogger() {
	logConfig := logger.Configuration{
		EnableConsole:     true,    // next, get from configuration
		ConsoleJSONFormat: true,    // next, get from configuration
		ConsoleLevel:      "debug", // next, get from configuration
	}

	if err := logger.NewLogger(logConfig, logger.InstanceZapLogger); err != nil {
		log.Fatalf("Could not instantiate log %v", err)
	}
}
