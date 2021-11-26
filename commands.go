package main

import (
	"errors"
	"fmt"
	"iftarbot/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ping command
func ping(session *discordgo.Session, event *discordgo.MessageCreate)(*discordgo.Message, error) {
	ping := session.HeartbeatLatency()
	pingString := strings.Replace(i18n.Lookup("en", "ping"), "{{ping}}", fmt.Sprintf("%v", ping.Milliseconds()), -1)
	return session.ChannelMessageSend(event.ChannelID, pingString)
}

// play command, plays adzan sound
func play(session *discordgo.Session, event *discordgo.MessageCreate)(*discordgo.Message, error) {
	println(isPlaying)
	if !isPlaying {
		if event.Author.ID != ownerID {
			return session.ChannelMessageSend(event.ChannelID, i18n.Lookup("en", "owner_only"))
		}

		err := playSound(session, guildID, channelID)
		if err != nil {
			format := fmt.Sprintf("An error occured, %v", err.Error())

			logger.WithFields(logger.Fields{"component": "commands", "action": "play adzan sound."}).
				Errorf(format)

			return session.ChannelMessageSend(event.ChannelID, format)
		}
	} 

	return nil, errors.New("audio is playing")
}

// iftar command, get iftar information
func iftar(session *discordgo.Session, event *discordgo.MessageCreate) (*discordgo.Message, error){
	var format string

	iftar, err := getIftar()
	if err != nil {
		format = fmt.Sprintf("An error occured while parsing date, %v", err.Error())

		logger.WithFields(logger.Fields{"component": "commands", "action": "get iftar information."}).
			Errorf(format)
	}

	iftarString := strings.Replace(i18n.Lookup("en", "iftar"), "{{time}}", fmt.Sprintf("%v:%v", iftar.Hour(), iftar.Minute()), -1)
	return session.ChannelMessageSend(event.ChannelID, iftarString)
}
