package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ready(session *discordgo.Session, event *discordgo.Ready) {
	session.UpdateGameStatus(0, "Running on GO! Ramadan special ✨");
}

func messageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID == session.State.User.ID || event.Author.Bot || !strings.HasPrefix(event.Content, prefix) {
		return;
	}

	if strings.HasPrefix(
		event.Content,
		fmt.Sprintf("%vping", prefix),
	) {
		ping_command(session, event);
	}

	if strings.HasPrefix(
		event.Content,
		fmt.Sprintf("%vplay", prefix),
	) {
		play_command(session, event);
	}

	if strings.HasPrefix(
		event.Content,
		fmt.Sprintf("%viftar", prefix),
	) {
		iftar_command(session, event);
	}
}
