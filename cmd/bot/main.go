package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/netscrawler/trap_bot/bot"
	"github.com/netscrawler/trap_bot/bot/handlers"
	"github.com/netscrawler/trap_bot/repository"
	"github.com/netscrawler/trap_bot/storage"
	"golang.org/x/net/proxy"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Error("BOT_TOKEN environment variable not set")
		return
	}

	opts := []telego.BotOption{}

	if botProxy := os.Getenv("BOT_PROXY"); botProxy != "" {
		client, err := buildHTTPClient(botProxy)
		if err != nil {
			log.ErrorContext(ctx, "create bot proxy", slog.Any("error", err))
			return
		}
		opts = append(opts, telego.WithHTTPClient(client))

	}

	if apiProxy := os.Getenv("API_SERVER"); apiProxy != "" {
		opts = append(opts, telego.WithAPIServer(apiProxy))
	}

	teleBot, err := telego.NewBot(botToken, opts...)
	if err != nil {
		log.Error("failed to create bot", "error", err)
		return
	}
	defer teleBot.Close(context.TODO())

	botUser, err := teleBot.GetMe(ctx)
	if err != nil {
		log.Error("failed to get bot user", "error", err)
		return
	}

	log.InfoContext(ctx, "started as", slog.Any("username", botUser.Username), slog.Any("id", botUser.ID))

	updates, err := teleBot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		log.Error("failed to get updates", "error", err)
		return
	}

	handler, err := th.NewBotHandler(teleBot, updates)
	if err != nil {
		log.Error("failed to create bot handler", "error", err)
		return
	}

	db, err := storage.NewStorage(ctx, "/app/data/db.db")
	if err != nil {
		log.Error("failed to create storage", "error", err)
		return
	}

	repo := repository.NewMediaFilesRepository(db)
	baseHandler := handlers.NewBaseHandler(teleBot, log)

	userHandler := handlers.NewUserHandler(baseHandler, repo)
	chatHandler := handlers.NewChatHandler(baseHandler, repo)
	router := bot.NewRouter(handler, userHandler, chatHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	router.Start()
	log.InfoContext(ctx, "router started")

	<-stop
	log.InfoContext(ctx, "stopping")

	if err := router.Stop(context.TODO()); err != nil {
		log.Error("failed to stop router", "error", err)
	}
}

func buildHTTPClient(proxyStr string) (*http.Client, error) {
	if proxyStr == "" {
		return &http.Client{}, nil
	}

	u, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {

	case "http", "https":
		return &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
			Timeout: 30 * time.Second,
		}, nil

	case "socks5", "socks5h":
		dialer, err := proxy.SOCKS5("tcp", u.Host, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}

		return &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
			},
			Timeout: 30 * time.Second,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s", u.Scheme)
	}
}
