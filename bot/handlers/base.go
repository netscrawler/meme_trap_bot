package handlers

import (
	"log/slog"

	"github.com/mymmrac/telego"
)

type BaseHandler struct {
	teleBot *telego.Bot
	log     *slog.Logger
}

func NewBaseHandler(teleBot *telego.Bot, log *slog.Logger) *BaseHandler {
	return &BaseHandler{
		teleBot: teleBot,
		log:     log,
	}
}
