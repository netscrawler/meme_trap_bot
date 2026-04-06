package bot

import (
	"context"

	th "github.com/mymmrac/telego/telegohandler"
	"github.com/netscrawler/trap_bot/bot/handlers"
)

type Router struct {
	h    *th.BotHandler
	user *handlers.UserHandler
	chat *handlers.ChatHandler
}

func NewRouter(
	botHandler *th.BotHandler,
	userHandler *handlers.UserHandler,
	chatHandler *handlers.ChatHandler,
) Router {
	base := botHandler.BaseGroup()
	base.Use(th.PanicRecovery())

	base.HandleInlineQuery(chatHandler.HandleInlineQuery)
	base.HandleMessage(userHandler.HandleIncomingMessage)

	// // -------------------- ADMIN --------------------
	// ah := base.Group(th.And(mw.AdminFilter,
	// 	th.Or(th.CommandEqual("admin"), th.CallbackDataPrefix("admin"))))

	// ah.HandleMessage(admin.Start, th.CommandEqual("admin"))
	// ah.HandleCallbackQuery(admin.MainMenu, th.CallbackDataEqual("admin_menu_start"))
	// ah.HandleCallbackQuery(admin.ListUser, th.CallbackDataEqual("admin_menu_list_users"))
	// ah.HandleCallbackQuery(admin.AddUser, th.CallbackDataEqual("admin_menu_add_user"))
	// ah.HandleCallbackQuery(admin.DeleteUser, th.CallbackDataEqual("admin_menu_delete_user"))

	// // Отдельная группа для текстовых сообщений админа
	// ahText := base.Group(mw.AdminStateFilter)
	// ahText.HandleMessage(admin.TextHandler, th.Not(th.AnyCommand()))

	// // -------------------- USER --------------------
	// uh := base.Group(th.And(
	// 	mw.UserFilter,
	// 	th.Or(th.CommandEqual("start"), th.CommandEqual("help")),
	// ))

	// uh.HandleMessage(user.Start, th.CommandEqual("start"))
	// uh.HandleMessage(user.Help, th.CommandEqual("help"))

	// inline := base.Group(mw.UserFilter)
	// inline.HandleInlineQuery(user.CheckStatusInlineQuery)

	// uhT := base.Group(mw.UserFilter)
	// uhT.HandleMessage(user.CheckStatusByAny, th.Not(th.AnyCommand()))

	return Router{
		h:    botHandler,
		user: userHandler,
		chat: chatHandler,
	}
}

func (r Router) Start() {
	if r.h.IsRunning() {
		return
	}

	go func() {
		_ = r.h.Start()
	}()
}

func (r Router) Stop(ctx context.Context) error {
	return r.h.StopWithContext(ctx)
}
