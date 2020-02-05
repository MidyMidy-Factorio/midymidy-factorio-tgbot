package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type playersResponse struct {
	Players []player `json:"players"`
}

type player struct {
	Name              string `json:"name"`
	Connected         bool   `json:"connected"`
	AfkTime           int    `json:"afk_time"`
	OnlineTime        int    `json:"online_time"`
	LastOnline        int    `json:"last_online"`
	DisplayResolution struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"display_resolution"`
	Spectator bool `json:"spectator"`
}

type updateResponse struct {
	Type        string `json:"type"`
	PlayerIndex int    `json:"player_index"`
	PlayerName  string `json:"player_name"`
	Message     string `json:"message"`
	Sarcasm     string `json:"sarcasm"`
}

type messageRequest struct {
	Messages []message `json:"messages"`
}

type message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func getUpdate() (*updateResponse, error) {
	resp, err := http.Get(webServiceBaseURL + "/update")
	if err != nil {
		return &updateResponse{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &updateResponse{}, err
	}

	ur := updateResponse{}
	err = json.Unmarshal(body, &ur)
	if err != nil {
		return &updateResponse{}, err
	}
	return &ur, nil
}

func getPlayers() (*[]player, error) {
	resp, err := http.Get(webServiceBaseURL + "/players")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pr := playersResponse{}
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return nil, err
	}
	return &pr.Players, nil
}

func getOnlinePlayers() (*[]player, error) {
	players, err := getPlayers()
	if err != nil {
		return nil, err
	}
	onlinePlayers := make([]player, 0)
	for _, p := range *players {
		if p.Connected {
			onlinePlayers = append(onlinePlayers, p)
		}
	}
	return &onlinePlayers, nil
}

func sendMessages(req messageRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(webServiceBaseURL+"/messages", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	/*
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	*/
	return nil
}

func sendMessage(name, msg string) error {
	m := message{name, msg}
	req := messageRequest{
		Messages: []message{m},
	}
	return sendMessages(req)
}
