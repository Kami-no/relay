package espace

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kami-no/relay/config"
)

type eSpace struct {
	Tenant  string `json:"tenant_id"`
	App     string `json:"app_id"`
	Theme   string `json:"theme_id"`
	Rcpt    string `json:"to_user_account"`
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"jump_url"`
	Type    int    `json:"type"`
}

type eResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func Send(rcpt string, title string, content string, url string) error {
	var espace eSpace

	espaceURL := config.Cfg.ESpace.URL
	espace.Tenant = config.Cfg.ESpace.Tenant
	espace.App = config.Cfg.ESpace.App
	espace.Theme = config.Cfg.ESpace.Theme

	espace.Rcpt = rcpt
	espace.Title = title
	espace.Content = content
	espace.URL = url

	espace.Type = 1

	requestBody, err := json.Marshal(espace)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		espaceURL,
		"application/json",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	var body eResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}

	if body.Code != "200" {
		log.Printf("Bad body %v", body)
	}

	defer resp.Body.Close()
	return nil
}
