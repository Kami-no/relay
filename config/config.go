package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

const cfgLocation = "config.yaml"

type Config struct {
	Direction string `json:"direction"`
	Gotify    string `json:"gotify"`
	ESpace    ESpace `json:"espace"`
}

type ESpace struct {
	URL    string `json:"url"`
	Tenant string `json:"tenant"`
	App    string `json:"app"`
	Theme  string `json:"theme"`
	Rcpt   string `json:"rcpt"`
}

func getConfig() *Config {
	var c Config

	yamlFile, err := ioutil.ReadFile(cfgLocation)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &c
}

var Cfg = getConfig()
