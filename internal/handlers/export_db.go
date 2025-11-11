package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/keyboards"
)

func ExportDBHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "📤 Експортую базу даних...",
	})

	stats, err := userRepo.GetStats(ctx)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("❌ Помилка отримання статистики: %v", err),
		})
		return
	}

	dbData, err := userRepo.ExportDB(ctx)
	if err != nil {
		log.Printf("Error exporting DB: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("❌ Помилка експорту БД: %v", err),
		})
		return
	}

	caption := fmt.Sprintf(
		"💾 <b>База даних бота</b>\n\n"+
			"📊 <b>Статистика:</b>\n"+
			"• Всього користувачів: %d\n"+
			"• Активних: %d\n"+
			"• Відписалися: %d\n"+
			"• Заблокували: %d\n\n"+
			"📦 Розмір: %.2f KB",
		stats.Total,
		stats.Active,
		stats.Unsubscribed,
		stats.Blocked,
		float64(len(dbData))/1024,
	)

	_, err = b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID: chatID,
		Document: &models.InputFileUpload{
			Filename: "bot.db",
			Data:     bytes.NewReader(dbData),
		},
		Caption:     caption,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboards.BackToAdminPanelKeyboard(),
	})

	if err != nil {
		log.Printf("Error sending DB file: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("❌ Помилка відправки файлу: %v", err),
		})
		return
	}

	log.Printf("✅ DB exported to admin %d (size: %d bytes)", update.Message.From.ID, len(dbData))
}
