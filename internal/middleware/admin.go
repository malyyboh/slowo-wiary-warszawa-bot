package middleware

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var adminIDs []int64

func InitAdmins() {
	adminIDsStr := os.Getenv("ADMIN_USER_IDS")
	if adminIDsStr == "" {
		log.Println("Warning: ADMIN_USER_IDS not set in .env")
		return
	}

	ids := strings.Split(adminIDsStr, ",")
	for _, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Warning: Invalid admin ID: %s", idStr)
			continue
		}
		adminIDs = append(adminIDs, id)
	}

	log.Printf("Loaded %d admin(s)", len(adminIDs))
}

func IsAdmin(userID int64) bool {
	for _, adminID := range adminIDs {
		if adminID == userID {
			return true
		}
	}
	return false
}

func AdminOnly(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var userID int64

		if update.Message != nil {
			userID = update.Message.From.ID
		} else if update.CallbackQuery != nil {
			userID = update.CallbackQuery.From.ID
		} else {
			return
		}

		if !IsAdmin(userID) {
			if update.Message != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "❌ У вас немає доступу до цієї команди.",
				})
			} else if update.CallbackQuery != nil {
				b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					Text:            "❌ У вас немає доступу до цієї функції.",
					ShowAlert:       true,
				})
			}
			return
		}

		next(ctx, b, update)
	}
}
