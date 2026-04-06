package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/netscrawler/trap_bot/domain"
)

type PhotoProvider interface {
	GetRandom(ctx context.Context) (domain.Image, error)
}

type ChatHandler struct {
	BaseHandler
	PhotoProvider PhotoProvider
}

func NewChatHandler(bh *BaseHandler, pp PhotoProvider) *ChatHandler {
	return &ChatHandler{
		BaseHandler:   *bh,
		PhotoProvider: pp,
	}
}

func (ch *ChatHandler) HandleInlineQuery(ctx *th.Context, query telego.InlineQuery) error {
	photo, err := ch.PhotoProvider.GetRandom(ctx.Context())
	if err != nil {
		return err
	}

	err = ctx.Bot().AnswerInlineQuery(ctx.Context(), &telego.AnswerInlineQueryParams{
		InlineQueryID: query.ID,
		Results: []telego.InlineQueryResult{
			&telego.InlineQueryResultCachedPhoto{
				Type:        "photo",
				ID:          rndID(),
				PhotoFileID: photo.FileID,
				Caption:     fmt.Sprintf("Добавил: @%s", photo.AddedBy),
			},
		},
		CacheTime: 0,
	})
	if err != nil {
		return err
	}

	return nil
}

func rnd() string {
	return fmt.Sprintf("%d", rand.Int63())
}

func rndID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
