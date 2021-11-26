package main

import (
	"encoding/json"
	"fmt"
	"iftarbot/logger"
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
	Aksam                  string `json:"Aksam"`
	AyinSekliURL           string `json:"AyinSekliURL"`
	Gunes                  string `json:"Gunes"`
	GunesBatis             string `json:"GunesBatis"`
	GunesDogus             string `json:"GunesDogus"`
	HicriTarihKisa         string `json:"HicriTarihKisa"`
	HicriTarihKisaIso8601  string `json:"HicriTarihKisaIso8601"`
	HicriTarihUzun         string `json:"HicriTarihUzun"`
	HicriTarihUzunIso8601  string `json:"HicriTarihUzunIso8601"`
	Ikindi                 string `json:"Ikindi"`
	Imsak                  string `json:"Imsak"`
	KibleSaati             string `json:"KibleSaati"`
	MiladiTarihKisa        string `json:"MiladiTarihKisa"`
	MiladiTarihKisaIso8601 string `json:"MiladiTarihKisaIso8601"`
	MiladiTarihUzun        string `json:"MiladiTarihUzun"`
	MiladiTarihUzunIso8601 string `json:"MiladiTarihUzunIso8601"`
	Ogle                   string `json:"Ogle"`
	Yatsi                  string `json:"Yatsi"`
}

// fetchData, fetching data from API
func fetchData() ([]Vakit, error) {
	res, err := http.Get(API_URL)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "utils", "action": "do task."}).
			Errorf("An error occured while fetching data, %v", err)

		return nil, err
	}

	// check if http code is 200 (OK)
	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			logger.WithFields(logger.Fields{"component": "utils", "action": "do task."}).
				Errorf("An error occured while reading data, %v", err)

			return nil, err
		}

		var data []Vakit
		err = json.Unmarshal(body, &data)
		if err != nil {
			logger.WithFields(logger.Fields{"component": "utils", "action": "do task."}).
				Errorf("An error occured while unmarshaling data, %v", err)

			return nil, err
		}

		return data, nil
	}

	return nil, fmt.Errorf("got http status code %d from server", res.StatusCode)
}

// playSound, plays adzan sound
func playSound(session *discordgo.Session, guildID string, channelID string) error {
	voiceConnection, err := session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	isPlaying = true
	dgvoice.PlayAudioFile(voiceConnection, "./assets/adzan.mp3", make(<-chan bool))

	voiceConnection.Close()
	err = voiceConnection.Disconnect()
	if err != nil {
		return err
	}

	isPlaying = false	// reset flag
	return nil
}

// getIftar, get iftar information
func getIftar() (time.Time, error) {
	date, err := time.Parse(time.RFC3339Nano, data[0].MiladiTarihUzunIso8601)
	if err != nil {
		return date, err
	}

	split := strings.Split(data[0].Aksam, ":")

	iftar, err := time.ParseDuration(fmt.Sprintf("%vh%vm", split[0], split[1]))
	if err != nil {
		return date, err
	}

	date = date.Add(iftar)
	return date, nil
}

// isIftar, validate iftar time
func isIftar() (bool, error) {
	date, err := getIftar()
	if err != nil {
		return false, err
	}

	return (date.Hour() == time.Now().Hour()) && (date.Minute() == time.Now().Minute()), nil
}


// setupJob, setting up job scheduler
func setupJob(session *discordgo.Session) {
	scheduler := gocron.NewScheduler(time.UTC)
	job, err := scheduler.Every(30).Seconds().Do(task, session)

	logger.WithFields(logger.Fields{"component": "utils", "action": "setup job."}).
		Infof("job : %v, error : %v", job, err)

	scheduler.StartAsync()
}

// task, execute task
func task(session *discordgo.Session) {
	now := time.Now()

	execute, err := isIftar()
	if err != nil {
		logger.WithFields(logger.Fields{"component": "utils", "action": "do task."}).
			Errorf("An error occured while checking iftar time on task, %v", err.Error())

		return
	}

	if execute {
		err := playSound(session, guildID, channelID)
		if err != nil && !isPlaying {
			logger.WithFields(logger.Fields{"component": "utils", "action": "do task."}).
				Errorf("An error occured while playing sound on task, %v", err.Error())

			return
		}
	}

	elapse := time.Since(now)
	logger.WithFields(logger.Fields{"component": "utils", "action": "do task."}).
		Infof("Task completed! Took %v ms ", elapse)
}
