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
	GetRandom(ctx context.Context) ([]domain.Image, error)
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

	var res []telego.InlineQueryResult

	for _, ph := range photo {
		res = append(res, &telego.InlineQueryResultCachedPhoto{
			Type:        "photo",
			ID:          rndID(),
			PhotoFileID: ph.FileID,
			Caption:     fmt.Sprintf("Добавил: @%s", ph.AddedBy),
		})
	}

	err = ctx.Bot().AnswerInlineQuery(ctx.Context(), &telego.AnswerInlineQueryParams{
		InlineQueryID: query.ID,
		Results:       res,
		CacheTime:     0,
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
