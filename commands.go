package main

import (
	"errors"
	"fmt"
	"iftarbot/logger"

	"github.com/bwmarrin/discordgo"
)

// ping command
func ping(session *discordgo.Session, event *discordgo.MessageCreate)(*discordgo.Message, error) {
	ping := session.HeartbeatLatency()
	format := fmt.Sprintf(":ping_pong: Pong! %vms", ping.Milliseconds())
	return session.ChannelMessageSend(event.ChannelID, format)
}

// play command, plays adzan sound
func play(session *discordgo.Session, event *discordgo.MessageCreate)(*discordgo.Message, error) {
	if !isPlaying {
		if event.Author.ID != ownerID {
			return session.ChannelMessageSend(event.ChannelID, "Bu komut sahiplere özeldir.")
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

	format = fmt.Sprintf("Kayseri merkez için iftar vakti, %v:%v", iftar.Hour(), iftar.Minute())
	return session.ChannelMessageSend(event.ChannelID, format)
}
