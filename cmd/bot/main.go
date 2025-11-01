package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
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
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "admin_", bot.MatchTypePrefix,
		middleware.AdminOnly(handlers.AdminCallbackHandler))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, handlers.CallbackHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.HelpHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/admin", bot.MatchTypeExact,
		middleware.AdminOnly(handlers.AdminPanelHandler))
	log.Println("Bot started...")
	b.Start(ctx)
}
