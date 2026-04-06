package handlers

import (
	"context"
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/netscrawler/trap_bot/domain"
)

type PhotoSaver interface {
	Save(ctx context.Context, imgs []domain.Image) error
}

type UserHandler struct {
	BaseHandler
	PhotoSaver PhotoSaver
}

func NewUserHandler(bh *BaseHandler, ps PhotoSaver) *UserHandler {
	return &UserHandler{
		BaseHandler: *bh,
		PhotoSaver:  ps,
	}
}

func (uh *UserHandler) HandleIncomingMessage(ctx *th.Context, message telego.Message) error {

	if message.Photo == nil {
		return nil
	}
	imgs := newPhoto(message.Photo, message.Chat.ID, message.From.Username)
	if err := uh.PhotoSaver.Save(ctx.Context(), imgs); err != nil {
		uh.log.Error("failed to save photo", slog.Any("err", err))

		_, err = uh.teleBot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: tu.ID(message.From.ID),
			Text:   "Failed to save photo. Please try again later.",
		})
		if err != nil {
			uh.log.Error("failed to send error message", slog.Any("err", err))
			return err
		}
	}

	_, err := uh.teleBot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: tu.ID(message.From.ID),
		Text:   "Photo saved successfully.",
	})
	if err != nil {
		uh.log.Error("failed to send success message", slog.Any("err", err))
		return err
	}

	return nil
}

func newPhoto(photo []telego.PhotoSize, chatID int64, chatName string) []domain.Image {
	if len(photo) == 0 {
		return nil
	}

	var max telego.PhotoSize
	for _, p := range photo {
		if p.Width*p.Height > max.Width*max.Height {
			max = p
		}
	}

	return []domain.Image{{
		AddedByID:    chatID,
		AddedBy:      chatName,
		FileID:       max.FileID,
		FileUniqueID: max.FileUniqueID,
		FileSize:     max.FileSize,
		Width:        max.Width,
		Height:       max.Height,
	}}
}
