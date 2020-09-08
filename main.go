package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Object struct {
	Kind       string `json:"object_kind"`
	Attributes struct {
		ID     int    `json:"iid"`
		Title  string `json:"title"`
		Action string `json:"action"`
	} `json:"object_attributes"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Gitlab-Event") == "Issue Hook" {
		var object Object
		_ = json.NewDecoder(r.Body).Decode(&object)
		w.WriteHeader(http.StatusOK)
		fmt.Printf("%v %v %v: %v\n", object.Kind, object.Attributes.ID, object.Attributes.Action, object.Attributes.Title)
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
