package chat

import (
	"context"
	"net/http"

	"github.com/Rocksus/pogo/configs"
)

type Repository interface {
	Init(config configs.ChatConfig) *Client
	GetHandler(ctx context.Context, chatbot *Client) func(w http.ResponseWriter, req *http.Request)
}
