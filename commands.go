package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ping_command(session *discordgo.Session, event *discordgo.MessageCreate) {
	ping := session.HeartbeatLatency()
	format := fmt.Sprintf(":ping_pong: Pong! %vms", ping.Milliseconds())
	session.ChannelMessageSend(event.ChannelID, format)
}

func play_command(session *discordgo.Session, event *discordgo.MessageCreate) {
	if is_playing {
		return;
	}
	if event.Author.ID != ownerID {
		session.ChannelMessageSend(event.ChannelID, "Bu komut sahiplere özeldir.");
		return;
	}
	err := play_sound(session, guildID, channelID);
	if err != nil {
		format := fmt.Sprintf("An error occured, %v", err.Error());
		session.ChannelMessageSend(event.ChannelID, format);
	}
}

func iftar_command(session *discordgo.Session, event *discordgo.MessageCreate) {
	iftar, err := get_iftar();
	if err != nil {
		format := fmt.Sprintf("An error occured while parsing date, %v", err.Error());
		session.ChannelMessageSend(event.ChannelID, format);
		return;
	}
	format := fmt.Sprintf("Kayseri merkez için iftar vakti, %v:%v", iftar.Hour(), iftar.Minute());
	session.ChannelMessageSend(event.ChannelID, format);
}
