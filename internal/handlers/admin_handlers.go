package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/conversation"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/keyboards"
	internalModels "github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
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

	log.Printf("AdminCallbackHandler: received callback '%s' from user %d", data, callback.From.ID)

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
		userID := callback.From.ID
		chatID := callback.Message.Message.Chat.ID
		StartAddEventDialog(ctx, b, userID, chatID)

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
		})
		return

	case "admin_delete_event":
		userID := callback.From.ID
		conv := conversation.GetManager()
		conv.SetState(userID, internalModels.StateAwaitingDeleteID)
		text = "🗑️ <b>Видалення події</b>\n\n" +
			"Введіть <b>ID події</b> для видалення:\n\n" +
			"Ви можете побачити ID в списку подій."
		keyboard = keyboards.BackToAdminPanelKeyboard()

	case "admin_confirm_delete":
		log.Println("Case: confirm_delete - calling handleDeleteConfirm")
		handleDeleteConfirm(ctx, b, callback)
		return

	case "admin_cancel_delete":
		log.Println("Case: cancel_delete - calling handleDeleteCancel")
		handleDeleteCancel(ctx, b, callback)
		return

	case "admin_users":
		text = getAdminUsersStatsText()
		keyboard = keyboards.AdminUsersKeyboard()

	case "admin_list_users":
		text = getAdminUsersListText()
		keyboard = keyboards.AdminUsersListKeyboard()

	default:
		log.Printf("Case: default - unknown command '%s'", data)
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
			formatEventDate(event.Date),
			event.ID,
		)
	}

	text += "\n💡 ✅ - опубліковано, 📝 - чернетка"

	return text
}

func DeleteEventHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	messageText := strings.TrimSpace(update.Message.Text)

	eventID, err := strconv.Atoi(messageText)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ Неправильний формат ID. Введіть число.",
		})
		return
	}

	event, err := eventRepo.GetByID(eventID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ Подію з таким ID не знайдено.",
		})
		return
	}

	conv := conversation.GetManager()
	conv.SetState(userID, internalModels.StateAwaitingDeleteConfirm)

	conv.GetConversation(userID).EventData.ID = eventID

	text := fmt.Sprintf(
		"🗑️ <b>Підтвердження видалення</b>\n\n"+
			"Ви дійсно хочете видалити цю подію?\n\n"+
			"<b>%s</b>\n"+
			"📅 %s\n"+
			"ID: %d",
		event.Title,
		formatEventDate(event.Date),
		event.ID,
	)

	keyboard := keyboards.DeleteConfirmKeyboard()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboard,
	})
}

func handleDeleteConfirm(ctx context.Context, b *bot.Bot, callback *models.CallbackQuery) {
	userID := callback.From.ID
	chatID := callback.Message.Message.Chat.ID

	conv := conversation.GetManager()
	conversation := conv.GetConversation(userID)

	if conversation == nil {
		log.Printf("Error: conversation is nil for user %d", userID)
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "❌ Помилка: дані втрачено",
			ShowAlert:       true,
		})
		return
	}

	eventID := conversation.EventData.ID
	log.Printf("Trying to delete event ID: %d", eventID)

	err := eventRepo.Delete(eventID)
	if err != nil {
		log.Printf("Error deleting event: %v", err)
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: callback.Message.Message.ID,
			Text:      "❌ Помилка видалення події.",
		})
		conv.ClearState(userID)
		return
	}

	log.Printf("Event %d deleted successfully", eventID)
	conv.ClearState(userID)

	text := fmt.Sprintf("✅ Подію (ID: %d) успішно видалено!", eventID)
	keyboard := keyboards.AdminPanelKeyboard()

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   callback.Message.Message.ID,
		Text:        text,
		ReplyMarkup: keyboard,
	})

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})
}

func handleDeleteCancel(ctx context.Context, b *bot.Bot, callback *models.CallbackQuery) {
	userID := callback.From.ID
	chatID := callback.Message.Message.Chat.ID

	conv := conversation.GetManager()
	conv.ClearState(userID)

	text := "❌ Видалення скасовано."
	keyboard := keyboards.AdminPanelKeyboard()

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   callback.Message.Message.ID,
		Text:        text,
		ReplyMarkup: keyboard,
	})

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})
}

func getAdminUsersStatsText() string {
	stats, err := userRepo.GetStats()
	if err != nil {
		log.Printf("Error getting user stats: %v", err)
		return "❌ Помилка отримання статистики користувачів"
	}

	text := "📊 <b>Статистика користувачів</b>\n\n"
	text += fmt.Sprintf("👥 Всього: <b>%d</b>\n", stats.Total)
	text += fmt.Sprintf("✅ Активних: <b>%d</b>\n", stats.Active)
	text += fmt.Sprintf("🔕 Відписалися: <b>%d</b>\n", stats.Unsubscribed)
	text += fmt.Sprintf("❌ Заблокували: <b>%d</b>\n", stats.Blocked)

	return text
}

func getAdminUsersListText() string {
	users, err := userRepo.GetAll()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return "❌ Помилка отримання списку користувачів"
	}

	if len(users) == 0 {
		return "📋 <b>Список користувачів</b>\n\nКористувачів поки що немає."
	}

	limit := 20
	total := len(users)

	text := "📋 <b>Список користувачів</b>\n\n"

	if total > limit {
		text += fmt.Sprintf("Показано перші %d з %d користувачів\n\n", limit, total)
	} else {
		text += fmt.Sprintf("Всього: <b>%d</b> користувачів\n\n", total)
	}

	for i, user := range users {
		if i >= limit {
			break
		}

		var status string
		if user.IsBlocked {
			status = "❌"
		} else if !user.IsActive {
			status = "🔕"
		} else {
			status = "✅"
		}

		username := user.Username
		if username == "" {
			username = "немає"
		} else {
			username = "@" + username
		}

		text += fmt.Sprintf(
			"%s <b>%d. %s</b> (%s)\n"+
				"    ID: %d | %s\n\n",
			status,
			i+1,
			user.FirstName,
			username,
			user.UserID,
			formatEventDate(user.SubscribedAt),
		)
	}

	text += "\n💡 ✅ - активний, 🔕 - відписався, ❌ - заблокував бота"

	return text
}
