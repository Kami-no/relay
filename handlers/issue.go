package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kami-no/relay/config"
	"github.com/Kami-no/relay/espace"
	"github.com/Kami-no/relay/gotify"
)

type Issue struct {
	Kind       string `json:"object_kind"`
	Attributes struct {
		ID     int    `json:"iid"`
		Title  string `json:"title"`
		Action string `json:"action"`
		URL    string `json:"url"`
	} `json:"object_attributes"`
}

func HIssue(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Gitlab-Event") == "Issue Hook" {
		var issue Issue

		_ = json.NewDecoder(r.Body).Decode(&issue)
		w.WriteHeader(http.StatusOK)

		log.Printf("%v %v %v: %v\n",
			issue.Kind,
			issue.Attributes.ID,
			issue.Attributes.Action,
			issue.Attributes.Title)

		title := fmt.Sprintf("%v %v %v",
			issue.Kind,
			issue.Attributes.ID,
			issue.Attributes.Action)

		message := fmt.Sprint(issue.Attributes.Title)

		if config.Cfg.Direction == "gotify" {
			var priority int
			if issue.Attributes.Action == "update" {
				priority = 0
			} else {
				priority = 5
			}

			gotify.Send(title, message, priority)

		} else {
			if issue.Attributes.Action == "open" {
				rcpt := config.Cfg.ESpace.Rcpt
				url := issue.Attributes.URL
				err := espace.Send(rcpt, title, message, url)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
