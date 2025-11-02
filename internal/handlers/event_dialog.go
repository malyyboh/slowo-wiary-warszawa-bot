package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/conversation"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/keyboards"
	internalModels "github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
)

func StartAddEventDialog(ctx context.Context, b *bot.Bot, userID int64, chatID int64) {
	conv := conversation.GetManager()
	conv.SetState(userID, internalModels.StateAwaitingTitle)

	text := "➕ <b>Додавання нової події</b>\n\n" +
		"Крок 1 з 6\n" +
		"Введіть <b>назву події</b>:\n\n" +
		"Для скасування натисніть /cancel"
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}

func HandleEventDialogMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	conv := conversation.GetManager()
	state := conv.GetState(userID)

	if text == "/cancel" {
		conv.ClearState(userID)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ Додавання події скасовано.",
		})
		return
	}

	switch state {
	case internalModels.StateAwaitingTitle:
		handleTitle(ctx, b, userID, chatID, text)
	case internalModels.StateAwaitingDate:
		handleDate(ctx, b, userID, chatID, text)
	case internalModels.StateAwaitingDesc:
		handleDescription(ctx, b, userID, chatID, text)
	case internalModels.StateAwaitingLocation:
		handleLocation(ctx, b, userID, chatID, text)
	case internalModels.StateAwaitingCategory:
		handleCategory(ctx, b, userID, chatID, text)
	case internalModels.StateAwaitingRegURL:
		handleRegistrationURL(ctx, b, userID, chatID, text)
	}
}

func handleTitle(ctx context.Context, b *bot.Bot, userID int64, chatID int64, title string) {
	conv := conversation.GetManager()
	conversation := conv.GetConversation(userID)
	conversation.EventData.Title = title

	conv.SetState(userID, internalModels.StateAwaitingDate)

	text := "✅ Назва збережена!\n\n" +
		"Крок 2 з 6\n" +
		"Введіть <b>дату та час події</b>:\n\n" +
		"Формат: <code>ДД.ММ.РРРР ГГ:ХХ</code> або <code>ДД.ММ.РРРР</code>\n" +
		"Приклади:\n" +
		"• <code>25.12.2025 16:00</code>\n" +
		"• <code>31.12.2025</code> (без часу)\n\n" +
		"Для скасування натисніть /cancel"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}

func handleDate(ctx context.Context, b *bot.Bot, userID int64, chatID int64, dateStr string) {
	dateStr = strings.TrimSpace(dateStr)

	eventDate, err := parseEventDate(dateStr)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text: "❌ Неправильний формат дати!\n\n" +
				"Використовуйте формат:\n" +
				"• <code>25.12.2025 16:00</code>\n" +
				"• <code>25.12.2025</code>\n\n" +
				"Спробуйте ще раз:",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	conv := conversation.GetManager()
	conversation := conv.GetConversation(userID)
	conversation.EventData.Date = eventDate

	conv.SetState(userID, internalModels.StateAwaitingDesc)

	text := "✅ Дата збережена!\n\n" +
		"Крок 3 з 6\n" +
		"Введіть <b>опис події</b>:\n\n" +
		"Для скасування натисніть /cancel"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}

func handleDescription(ctx context.Context, b *bot.Bot, userID int64, chatID int64, description string) {
	conv := conversation.GetManager()
	conversation := conv.GetConversation(userID)
	conversation.EventData.Description = description

	conv.SetState(userID, internalModels.StateAwaitingLocation)

	text := "✅ Опис збережено!\n\n" +
		"Крок 4 з 6\n" +
		"Введіть <b>місце проведення</b> (адресу):\n\n" +
		"Або натисніть /skip щоб пропустити\n" +
		"Для скасування натисніть /cancel"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}

func handleLocation(ctx context.Context, b *bot.Bot, userID int64, chatID int64, location string) {
	conv := conversation.GetManager()
	conversation := conv.GetConversation(userID)

	if location != "/skip" {
		conversation.EventData.Location = &location
	}

	conv.SetState(userID, internalModels.StateAwaitingCategory)

	text := "✅ Місце збережено!\n\n" +
		"Крок 5 з 6\n" +
		"Введіть <b>категорію події</b>:\n\n" +
		"Наприклад: Богослужіння, Семінар, Концерт, Молодіжка\n\n" +
		"Або натисніть /skip щоб пропустити\n" +
		"Для скасування натисніть /cancel"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}

func handleCategory(ctx context.Context, b *bot.Bot, userID int64, chatID int64, category string) {
	conv := conversation.GetManager()
	conversation := conv.GetConversation(userID)

	if category != "/skip" {
		conversation.EventData.Category = &category
	}

	conv.SetState(userID, internalModels.StateAwaitingRegURL)

	text := "✅ Категорія збережена!\n\n" +
		"Крок 6 з 6\n" +
		"Введіть <b>посилання для реєстрації</b>:\n\n" +
		"Наприклад: https://forms.google.com/...\n\n" +
		"Або натисніть /skip щоб пропустити\n" +
		"Для скасування натисніть /cancel"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}

func handleRegistrationURL(ctx context.Context, b *bot.Bot, userID int64, chatID int64, url string) {
	conv := conversation.GetManager()
	conversation := conv.GetConversation(userID)

	if url != "/skip" {
		conversation.EventData.RegistrationURL = &url
	}

	conversation.EventData.IsPublished = true
	conversation.EventData.CreatedAt = time.Now()
	conversation.EventData.CreatedBy = userID

	err := eventRepo.Create(conversation.EventData)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ Помилка збереження події в базу даних.",
		})
		conv.ClearState(userID)
		return
	}

	summary := formatEventSummary(conversation.EventData)

	conv.ClearState(userID)

	keyboard := keyboards.AdminPanelKeyboard()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        summary,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboard,
	})
}

func parseEventDate(input string) (time.Time, error) {
	t, err := time.Parse("02.01.2006 15:04", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("02.01.2006", input)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("неправильний формат дати")
}

func formatEventSummary(event *internalModels.Event) string {
	text := "✅ <b>Подію успішно створено!</b>\n\n" +
		fmt.Sprintf("<b>%s</b>\n", event.Title) +
		fmt.Sprintf("📅 %s\n", formatEventDate(event.Date)) +
		fmt.Sprintf("📝 %s\n", event.Description)

	if event.Location != nil && *event.Location != "" {
		text += fmt.Sprintf("📍 %s\n", *event.Location)
	}
	if event.Category != nil && *event.Category != "" {
		text += fmt.Sprintf("🏷 %s\n", *event.Category)
	}
	if event.RegistrationURL != nil && *event.RegistrationURL != "" {
		text += fmt.Sprintf("🔗 %s\n", *event.RegistrationURL)
	}

	text += fmt.Sprintf("\nID події: %d", event.ID)

	return text
}
