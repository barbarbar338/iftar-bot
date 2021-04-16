package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	token string
	prefix string
	guildID string
	channelID string
	data []Vakit
)

func main() {
	res, err := fetch_data();
	if err != nil {
		fmt.Printf(
			"An error occured while loading dataset, %v \n", 
			err.Error(),
		);
		return;
	}
	data = res;

	err = godotenv.Load();
	if err != nil {
		fmt.Printf(
			"An error occured while loading .env file, %v \n", 
			err.Error(),
		);
		return;
	}

	token = os.Getenv("BOT_TOKEN");
	prefix = os.Getenv("BOT_PREFIX");
	guildID = os.Getenv("GUILD_ID");
	channelID = os.Getenv("CHANNEL_ID");

	discord, err := discordgo.New(fmt.Sprintf("Bot %v", token));
	if err != nil {
		fmt.Printf(
			"An error occured while creatind Discord session, %v \n", 
			err.Error(),
		);
		return;
	}

	discord.AddHandler(ready);
	discord.AddHandler(messageCreate);

	err = discord.Open();

	if err != nil {
		fmt.Printf(
			"An error occured while connectin to Discord API, %v \n",
			err.Error(),
		);
		return;
	}

	fmt.Printf(
		"Logged in as %v \n",
		discord.State.User.Username,
	);

	setup_job(discord);

	sc := make(chan os.Signal, 1);
	signal.Notify(
		sc, 
		syscall.SIGINT, 
		syscall.SIGTERM, 
		os.Interrupt,
	);
	<-sc;

	fmt.Println("Closing bot, see you later...");

	discord.Close();
}


