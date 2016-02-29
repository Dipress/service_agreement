package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Configuration struct {
	DBUsername  string
	DBPassword  string
	DBName      string
	DBDriver    string
	DBProtocol  string
	DBHost      string
	DBPort      string
	SSHUser     string
	SSHPassword string
	SSHHost     string
	SSHPort     string
	SSHProtocol string
	FileName    string
	RemotePath  string
}

var cfgInstance *Configuration
var cfgOnce sync.Once

func GetConfiguration() *Configuration {
	cfgOnce.Do(func() {
		file, _ := os.Open("config.json")
		decoder := json.NewDecoder(file)
		cfg := &Configuration{}
		err := decoder.Decode(&cfg)

		if err != nil {
			log.Println("Config error")
			log.Fatal(err)
		}
		cfgInstance = cfg
	})
	return cfgInstance
}
