package gotify

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kami-no/relay/config"
)

type Gotify struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

func Send(title string, message string, priority int) {
	var gotify Gotify

	gotify.Title = title
	gotify.Message = message
	gotify.Priority = priority

	requestBody, err := json.Marshal(gotify)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.Post(
		config.Cfg.Gotify,
		"application/json",
		bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
