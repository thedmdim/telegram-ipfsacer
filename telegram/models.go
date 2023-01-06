package telegram

import "encoding/json"

type Updates struct {
	Results []Result `json:"result"`
}

type Result struct {
	UpdateId int     `json:"update_id"`
	Post  Post `json:"channel_post"`
}

type Post struct {
	Id json.Number `json:"message_id"`
	Text           string   `json:"text"`
}

type EditedPost struct {
	ChatId string `json:"chat_id"`
	Text   string      `json:"text"`
	Id json.Number `json:"message_id"`
	ParseMode string `json:"parse_mode"`
}