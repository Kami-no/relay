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
	Kind string `json:"object_kind"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var object Object
	_ = json.NewDecoder(r.Body).Decode(&object)
	w.WriteHeader(http.StatusOK)
	fmt.Printf("%s \"%s\" %s\n", r.RemoteAddr, r.Header.Get("X-Gitlab-Event"), object.Kind)
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
