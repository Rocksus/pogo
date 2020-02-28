package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Rocksus/pogo/internal/modules/joke"
	"github.com/Rocksus/pogo/internal/modules/news"
	"github.com/Rocksus/pogo/internal/modules/weather"
	"github.com/Rocksus/pogo/internal/repositories/chat"
	"github.com/Rocksus/pogo/internal/repositories/interpretor"
	"github.com/Rocksus/pogo/internal/utils/logging"

	"github.com/Rocksus/pogo/configs"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("[init cmd] environment file not found.")
	} else {
		log.Print("[init cmd] successfully loaded environment files.")
	}
}

func main() {
	config := configs.New()

	weather.Init(config.Weather)
	news.Init(config.News)
	joke.Init()

	interpretor := interpretor.InitInterpretorRepository(config.Interpretor)

	chatbot := chat.InitChatRepository(config.Chat, interpretor)
	handler := logging.Middleware(chatbot.GetHandler())

	http.HandleFunc("/callback", handler)

	srv := &http.Server{
		Handler:      http.DefaultServeMux,
		Addr:         ":" + config.Port,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Listening on port %s\n", config.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("[Init]Fail to start serving port %s, err: %v", config.Port, err)
	}
}
