package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/keyboards"
)

func AdminPanelHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := "🔐 <b>Адмін-панель</b>\n\nОберіть дію:"
	keyboard := keyboards.AdminPanelKeyboard()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboard,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: bot.True(),
		},
	})
}

func AdminCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	data := callback.Data

	var text string
	var keyboard *models.InlineKeyboardMarkup

	switch data {
	case "admin_panel":
		text = "🔐 <b>Адмін-панель</b>\n\nОберіть дію:"
		keyboard = keyboards.AdminPanelKeyboard()

	case "admin_list_events":
		text = getAdminEventsListText()
		keyboard = keyboards.AdminEventsListKeyboard()

	case "admin_add_event":
		text = "➕ <b>Додавання нової події</b>\n\n" +
			"Ця функція буде реалізована в наступному кроці.\n" +
			"Для додавання події буде використано діалог з ботом."
		keyboard = keyboards.BackToAdminPanelKeyboard()

	default:
		text = "Невідома команда"
		keyboard = keyboards.AdminPanelKeyboard()
	}

	if callback.Message.Message == nil {
		log.Printf("Error: callback message is nil")
		return
	}

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      callback.Message.Message.Chat.ID,
		MessageID:   callback.Message.Message.ID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboard,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: bot.True(),
		},
	})

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})
}

func getAdminEventsListText() string {
	events, err := eventRepo.GetAll()
	if err != nil {
		log.Printf("Error getting events: %v", err)
		return "❌ Помилка отримання подій з бази даних"
	}

	if len(events) == 0 {
		return "📋 <b>Список подій</b>\n\nПодій поки що немає."
	}

	text := "📋 <b>Список подій</b>\n\n"

	for i, event := range events {
		status := "✅"
		if !event.IsPublished {
			status = "📝"
		}

		text += fmt.Sprintf(
			"%s <b>%d. %s</b>\n"+
				"📅 %s\n"+
				"ID: %d\n\n",
			status,
			i+1,
			event.Title,
			event.Date.Format("02.01.2006 15:04"),
			event.ID,
		)
	}

	text += "\n💡 ✅ - опубліковано, 📝 - чернетка"

	return text
}
