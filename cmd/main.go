package main

import (
	"log"
	"net/http"

	"github.com/Rocksus/pogo/configs"
	"google.golang.org/api/chat/v1"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("[init cmd] environment file not found.")
	} else {
		log.Print("[init cmd] successfully loaded environment files.")
	}
}

func main() {
	config := configs.New()

	chatbot, err := chat.Init(config.Chat)
	if err != nil {
		log.Fatalf("[Init]Fail to initialize chat repository, err: %v", err)
	}
	http.HandleFunc("/callback", chat.GetHandler(chatbot))

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("[Init]Fail to start serving port %s, err: %v", config.Port, err)
	}
}
