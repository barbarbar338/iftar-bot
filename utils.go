package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
)

const (
	API_URL = "https://ezanvakti.herokuapp.com/vakitler/9620" // Kayseri, merkez.
)

type Vakit struct {
	Aksam string `json:"Aksam"`
    AyinSekliURL string `json:"AyinSekliURL"`
    Gunes string `json:"Gunes"`
    GunesBatis string `json:"GunesBatis"`
    GunesDogus string `json:"GunesDogus"`
    HicriTarihKisa string `json:"HicriTarihKisa"`
    HicriTarihKisaIso8601 string `json:"HicriTarihKisaIso8601"`
    HicriTarihUzun string `json:"HicriTarihUzun"`
    HicriTarihUzunIso8601 string `json:"HicriTarihUzunIso8601"`
    Ikindi string `json:"Ikindi"`
    Imsak string `json:"Imsak"`
    KibleSaati string `json:"KibleSaati"`
    MiladiTarihKisa string `json:"MiladiTarihKisa"`
    MiladiTarihKisaIso8601 string `json:"MiladiTarihKisaIso8601"`
    MiladiTarihUzun string `json:"MiladiTarihUzun"`
    MiladiTarihUzunIso8601 string `json:"MiladiTarihUzunIso8601"`
    Ogle string `json:"Ogle"`
    Yatsi string `json:"Yatsi"`
}

func fetch_data() ([]Vakit, error) {
	res, err := http.Get(API_URL);
	if err != nil {
		fmt.Printf("An error occured while fetching data, %v \n", err.Error());
		return nil, err;
	}

	body, err := io.ReadAll(res.Body);
	if err != nil {
		fmt.Printf("An error occured while reading data, %v \n", err.Error());
		return nil, err;
	}

	var data []Vakit;
	
	err = json.Unmarshal(body, &data);
	if err != nil {
		fmt.Printf("An error occured while unmarshaling data, %v \n", err.Error());
		return nil, err;
	}

	return data, nil;
}

func play_sound(session *discordgo.Session, guildID string, channelID string) error {
	voice_connection, err := session.ChannelVoiceJoin(
		guildID, 
		channelID, 
		false, 
		true,
	);
	if err != nil {
		return err;
	}
	is_playing = true;
	dgvoice.PlayAudioFile(
		voice_connection,
		"./assets/ezan.mp3",
		make(<-chan bool),
	);
	voice_connection.Close();
	err = voice_connection.Disconnect();
	if err != nil {
		return err;
	}
	is_playing = false;
	return nil;
}

func get_iftar() (time.Time, error) {
	date, err := time.Parse(time.RFC3339Nano, data[0].MiladiTarihUzunIso8601);
	if err != nil {
		return date, err;
	}
	
	split := strings.Split(data[0].Aksam, ":");

	iftar, err := time.ParseDuration(fmt.Sprintf("%vh%vm", split[0], split[1]));
	if err != nil {
		return date, err;
	}

	date = date.Add(iftar);

	return date, nil;
}

func is_iftar() (bool, error) {
	date, err := get_iftar();
	if err != nil {
		return false, err;
	}

	now := time.Now();

	return (date.Hour() == now.Hour()) && (date.Minute() == now.Minute()), nil;
}

func setup_job(session *discordgo.Session) {
	scheduler := gocron.NewScheduler(time.UTC);
	scheduler.Every(30).Seconds().Do(task, session);
	scheduler.StartAsync();
}

func task(session *discordgo.Session) {
	now := time.Now();

	execute, err := is_iftar();
	if err != nil {
		fmt.Printf("An error occured while checking iftar time on task, %v \n", err.Error());
		return;
	}

	if execute {
		err := play_sound(session, guildID, channelID);
		if err != nil && !is_playing {
			fmt.Printf("An error occured while playing sound on task, %v \n", err.Error());
			return;
		}
	}

	then := time.Now();
	fmt.Printf("Task completed! Took %vms \n", then.Unix() - now.Unix());
}
