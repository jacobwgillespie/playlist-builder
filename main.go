package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/shkh/lastfm-go/lastfm"
)

var router = mux.NewRouter()
var DB gorm.DB
var pool *redis.Pool
var api *lastfm.Api

var buildJobs = make(chan BuildArguments)

func main() {
	// Read config from flags
	flag.Int64Var(&Config.Port, "port", 8080, "port to listen on")
	flag.StringVar(&Config.Dsn, "dsn", "root@tcp(127.0.0.1:3306)/playlist?charset=utf8&parseTime=True", "the database DSN (like root@tcp(127.0.0.1:3306)/playlist?charset=utf8&parseTime=True)")
	flag.StringVar(&Config.RedisServer, "redisServer", ":6379", "")
	flag.StringVar(&Config.RedisPassword, "redisPassword", "", "")
	flag.StringVar(&Config.LastFMAPIKey, "lastfmApiKey", "", "Last.FM API key")
	flag.Parse()

	// Read config from ENV
	if port := os.Getenv("PORT"); port != "" {
		Config.Port, _ = strconv.ParseInt(port, 10, 64)
	}
	if dsn := os.Getenv("DSN"); dsn != "" {
		Config.Dsn = dsn
	}
	if redisServer := os.Getenv("REDIS_SERVER"); redisServer != "" {
		Config.RedisServer = redisServer
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		Config.RedisPassword = redisPassword
	}
	if lastfmKey := os.Getenv("LASTFM_API_KEY"); lastfmKey != "" {
		Config.LastFMAPIKey = lastfmKey
	}

	// Build the workers
	wg := new(sync.WaitGroup)
	for w := 1; w <= 50; w++ {
		wg.Add(1)
		go worker(w, buildJobs, wg)
	}

	api = lastfm.New(Config.LastFMAPIKey, "")

	// Connect to the database
	var err error
	DB, err = gorm.Open("mysql", Config.Dsn)
	if err != nil {
		log.Panicf("ERROR: %v\n", err)
	}
	DB.LogMode(true)
	pool = newRedisPool(Config.RedisServer, Config.RedisPassword)

	router.Handle("/api/build", requestHandler(handleApiBuild))
	// Build the middleware chain
	chain := methodOverrideMiddleware(corsMiddleware(router))

	// Start the server
	http.Handle("/", chain)
	log.Printf("Server listening on %d\n", Config.Port)
	http.ListenAndServe(":"+strconv.FormatInt(Config.Port, 10), nil)
}
