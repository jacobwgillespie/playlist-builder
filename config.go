package main

type configStruct struct {
	Port          int64
	Dsn           string
	RedisServer   string
	RedisPassword string
	LastFMAPIKey  string
}

var Config configStruct
