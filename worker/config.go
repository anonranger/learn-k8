package main

import "os"

type Config struct {
	MongoURI string
	PORT     string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		MongoURI: os.Getenv("MONGO_URI"),
		PORT:     os.Getenv("PORT"),
	}
}
