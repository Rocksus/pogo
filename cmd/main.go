package main

import (
	"log"
	"net/http"

	"github.com/Rocksus/pogo/internal/repositories/chat"
	"github.com/Rocksus/pogo/internal/repositories/interpretor"

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

	interpretor := interpretor.InitInterpretorRepository(config.Interpretor)

	chatbot := chat.InitChatRepository(config.Chat, interpretor)

	http.HandleFunc("/callback", chatbot.GetHandler())

	log.Printf("Listening on port %s\n", config.Port)
	if err := http.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile, nil); err != nil {
		log.Fatalf("[Init]Fail to start serving port %s, err: %v", config.Port, err)
	}
}
