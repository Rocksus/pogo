package main

import (
	"net/http"
	"time"

	"github.com/Rocksus/pogo/internal/controllers/linehttp"
	"github.com/Rocksus/pogo/internal/modules/joke"
	"github.com/Rocksus/pogo/internal/modules/news"
	"github.com/Rocksus/pogo/internal/modules/weather"
	"github.com/Rocksus/pogo/internal/repositories/interpreter/witai"
	"github.com/Rocksus/pogo/internal/usecase/replier"
	"github.com/Rocksus/pogo/internal/utils/logging"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"

	"github.com/Rocksus/pogo/configs"
	"github.com/joho/godotenv"
)

func main() {
	log.Init(log.WithLevel(log.DebugLevel))

	if err := godotenv.Load(); err != nil {
		log.Fatalln("environment file not found")
	}

	log.Infoln("successfully loaded environment files")

	config := configs.New()

	weather.Init(config.Weather)
	news.Init(config.News)
	joke.Init()

	bot, err := linebot.New(config.Chat.ChannelSecret, config.Chat.ChannelAccessToken)
	if err != nil {
		log.WithError(err).Fatalln("failed to create linebot client")
	}

	interpreter := witai.NewInterpreter(config.Interpretor)
	replier := replier.NewMessageReplier(interpreter)
	controller := linehttp.NewController(bot, replier)
	http.HandleFunc("/callback", logging.Middleware(controller.HandleWebhook))

	srv := &http.Server{
		Handler:      http.DefaultServeMux,
		Addr:         ":" + config.Port,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Infof("Listening on port %s", config.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatalf("Failed to start serving port %s", config.Port)
	}
}
