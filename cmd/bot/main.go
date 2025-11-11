package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/conversation"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/database"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/handlers"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/middleware"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	middleware.InitAdmins()
	conversation.InitManager()

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/bot.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
		bot.WithMiddlewares(
			middleware.RateLimit(2*time.Second, 20, 1*time.Minute),
		),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	_, err = b.DeleteWebhook(ctx, &bot.DeleteWebhookParams{
		DropPendingUpdates: true,
	})
	if err != nil {
		log.Printf("Failed to delete webhook: %v", err)
	} else {
		log.Printf("✅ Webhook deleted, old updates dropped")
	}

	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Головне меню"},
			{Command: "help", Description: "Довідка бота"},
			{Command: "menu", Description: "Показати кнопки меню"},
			{Command: "privacy", Description: "Політика конфіденційності"},
		},
	})
	if err != nil {
		log.Printf("Failed to set bot commands: %v", err)
	}

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "admin_", bot.MatchTypePrefix,
		middleware.AdminOnly(handlers.AdminCallbackHandler))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, handlers.CallbackHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.HelpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/menu", bot.MatchTypeExact, handlers.MenuHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/privacy", bot.MatchTypeExact, handlers.PrivacyHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/admin", bot.MatchTypeExact,
		middleware.AdminOnly(handlers.AdminPanelHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/exportDB", bot.MatchTypeExact, middleware.AdminOnly(handlers.ExportDBHandler))

	log.Println("Bot started...")
	b.Start(ctx)
}
