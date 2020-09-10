package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gopkg.in/yaml.v3"
)

type Issue struct {
	Kind       string `json:"object_kind"`
	Attributes struct {
		ID     int    `json:"iid"`
		Title  string `json:"title"`
		Action string `json:"action"`
	} `json:"object_attributes"`
}

type Gotify struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

type config struct {
	Endpoint string `json:"endpoint"`
}

func (c *config) getConfig() *config {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Gitlab-Event") == "Issue Hook" {
		var issue Issue
		_ = json.NewDecoder(r.Body).Decode(&issue)
		w.WriteHeader(http.StatusOK)

		log.Printf("%v %v %v: %v\n",
			issue.Kind,
			issue.Attributes.ID,
			issue.Attributes.Action,
			issue.Attributes.Title)

		var gotify Gotify
		gotify.Title = fmt.Sprintf("%v %v %v",
			issue.Kind,
			issue.Attributes.ID,
			issue.Attributes.Action)

		gotify.Message = fmt.Sprint(issue.Attributes.Title)

		if issue.Attributes.Action == "update" {
			gotify.Priority = 0
		} else {
			gotify.Priority = 5
		}

		requestBody, err := json.Marshal(gotify)
		if err != nil {
			log.Println(err)
			return
		}

		var cfg config
		cfg.getConfig()

		resp, err := http.Post(
			cfg.Endpoint,
			"application/json",
			bytes.NewBuffer(requestBody))
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
	}
}

func main() {
	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
