package main

import (
	"fmt"
	"iftarbot/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ready event, update status
func ready(session *discordgo.Session, event *discordgo.Ready) {
	err := session.UpdateGameStatus(0, "Running on GO! Ramadan special ✨")
	if err != nil {
		logger.WithFields(logger.Fields{"component": "events", "action": "update status to ready."}).
			Errorf("An error occurred while update status to ready, %v", err)
	}
}

// messageCreate event, handle commands
func messageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID == session.State.User.ID || event.Author.Bot || !strings.HasPrefix(event.Content, prefix) {
		return
	}

	if strings.HasPrefix(event.Content, fmt.Sprintf("%vping", prefix)) {
		_, err := ping(session, event)
		if err != nil {
			logger.WithFields(logger.Fields{"component": "events", "action": "send ping message."}).
				Errorf("An error occurred while sending ping message, %v", err)
		}
	}

	if strings.HasPrefix(event.Content, fmt.Sprintf("%vplay", prefix)) {
		_, err := play(session, event)
		if err != nil {
			logger.WithFields(logger.Fields{"component": "events", "action": "play adzan sound."}).
				Errorf("An error occurred while play adzan sound, %v", err)
		}
	}

	if strings.HasPrefix(event.Content, fmt.Sprintf("%viftar", prefix)) {
		_, err := iftar(session, event)
		if err != nil {
			logger.WithFields(logger.Fields{"component": "events", "action": "get iftar information."}).
				Errorf("An error occurred when get iftar information, %v", err)
		}
	}
}
