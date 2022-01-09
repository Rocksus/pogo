package main

import (
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/nickylogan/go-log"

	"github.com/Rocksus/pogo/internal/controllers/linehttp"
	"github.com/Rocksus/pogo/internal/repositories/interpreter/witai"
	"github.com/Rocksus/pogo/internal/usecase/gauth"
	"github.com/Rocksus/pogo/internal/usecase/replier"
	"github.com/Rocksus/pogo/internal/utils/logging"
	"github.com/Rocksus/pogo/pkg/plugin"
	"github.com/Rocksus/pogo/pkg/plugin/joke"
	"github.com/Rocksus/pogo/pkg/plugin/moneysheets"
	"github.com/Rocksus/pogo/pkg/plugin/weather"

	"github.com/joho/godotenv"

	"github.com/Rocksus/pogo/configs"
)

func main() {
	log.Init(log.WithLevel(log.DebugLevel))

	if err := godotenv.Load(); err != nil {
		log.Fatalln("environment file not found")
	}

	log.Infoln("successfully loaded environment files")

	config := configs.New()

	weatherPlugin := weather.InitPlugin(config.Weather.APIKey)
	jokePlugin := joke.InitPlugin()

	gauth, err := gauth.New(config.Google,
		moneysheets.GetScopes)
	if err != nil {
		log.WithError(err).Fatalln("failed to create gauth client")
	}
	moneySheets, err := moneysheets.InitPlugin(config.MoneySheets, gauth.GetClient())
	if err != nil {
		log.WithError(err).Fatalln("failed to create moneySheets plugin")
	}

	bot, err := linebot.New(config.Chat.ChannelSecret, config.Chat.ChannelAccessToken)
	if err != nil {
		log.WithError(err).Fatalln("failed to create linebot client")
	}

	interpreter := witai.NewInterpreter(config.Interpretor)
	replier := replier.NewMessageReplier(interpreter, map[string]plugin.MessageReplier{
		"tellJoke":                          jokePlugin,
		"weather/checkWeather":              weatherPlugin,
		"financialPlanning_addBalance":      moneySheets.GetTransactionAdder(),
		"financialPlanning/checkBalance":    moneySheets.GetBalanceChecker(),
		"financialPlanning/insertSpendings": moneySheets.GetSpendingAdder(),
	})
	controller := linehttp.NewController(bot, replier)
	http.HandleFunc("/callback", logging.Middleware(controller.HandleWebhook))

	srv := &http.Server{
		Handler:      http.DefaultServeMux,
		Addr:         "0.0.0.0:" + config.Port,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Infof("Listening on port %s", config.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatalf("Failed to start serving port %s", config.Port)
	}
}
