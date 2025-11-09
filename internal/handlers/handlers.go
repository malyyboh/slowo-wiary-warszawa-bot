package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/conversation"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/keyboards"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/messages"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/middleware"
	internalModels "github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/repository"
)

var eventRepo = repository.NewEventRepository()
var userRepo = repository.NewUserRepository()

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	username := update.Message.From.Username
	firstName := update.Message.From.FirstName

	location, _ := time.LoadLocation("Europe/Warsaw")
	now := time.Now().In(location)

	user := &internalModels.User{
		UserID:       userID,
		Username:     username,
		FirstName:    firstName,
		SubscribedAt: now,
		IsActive:     true,
		IsBlocked:    false,
		LastSeen:     now,
		UpdatedAt:    now,
	}

	err := userRepo.AddOrUpdate(ctx, user)
	if err != nil {
		log.Printf("Error adding/updating user: %v", err)
	}

	isActive := true
	savedUser, err := userRepo.GetByID(ctx, userID)
	if err == nil && savedUser != nil {
		isActive = savedUser.IsActive
	}

	text := messages.GetText("/start")
	keyboard := keyboards.MainMenuReplyKeyboard(isActive)

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

func HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := messages.GetText("/help")
	keyboard := keyboards.BackToMainMenuKeyboard()

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

func MenuHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	isActive := true
	user, err := userRepo.GetByID(ctx, userID)
	if err == nil && user != nil {
		isActive = user.IsActive
	}

	text := "📱 <b>Головне меню</b>\n\nОберіть розділ з кнопок нижче:"
	keyboard := keyboards.MainMenuReplyKeyboard(isActive)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboard,
	})
}

func CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	data := callback.Data

	var text string
	var keyboard *models.InlineKeyboardMarkup

	switch data {
	case "back_to_start":
		text = messages.GetText("/start")
		keyboard = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{},
		}

	case "about_us":
		text = messages.GetText("about_menu")
		keyboard = keyboards.AboutUsKeyboard()
	case "about_church":
		text = messages.GetText("about_church")
		keyboard = keyboards.BackToAboutKeyboard()
	case "about_church_mission":
		text = messages.GetText("about_church_mission")
		keyboard = keyboards.BackToAboutKeyboard()
	case "about_church_doctrine":
		text = messages.GetText("about_church_doctrine")
		keyboard = keyboards.BackToAboutKeyboard()
	case "about_church_pastors":
		text = messages.GetText("about_church_pastors")
		keyboard = keyboards.BackToAboutKeyboard()
	case "about_church_history":
		text = messages.GetText("about_church_history")
		keyboard = keyboards.BackToAboutKeyboard()

	case "ministry":
		text = messages.GetText("ministry_menu")
		keyboard = keyboards.MinistryKeyboard()
	case "sunday_ministry":
		text = messages.GetText("sunday_ministry")
		keyboard = keyboards.BackToMinistryKeyboard()
	case "home_ministry":
		text = messages.GetText("home_ministry")
		keyboard = keyboards.BackToMinistryKeyboard()
	case "prayer_ministry":
		text = messages.GetText("prayer_ministry")
		keyboard = keyboards.BackToMinistryKeyboard()
	case "youth_ministry":
		text = messages.GetText("youth_ministry")
		keyboard = keyboards.BackToMinistryKeyboard()
	case "teenagers_ministry":
		text = messages.GetText("teenagers_ministry")
		keyboard = keyboards.BackToMinistryKeyboard()
	case "kindergarten_ministry":
		text = messages.GetText("kindergarten_ministry")
		keyboard = keyboards.BackToMinistryKeyboard()

	case "social_media":
		text = messages.GetText("social_media")
		keyboard = keyboards.BackToMainMenuKeyboard()
	case "donation":
		text = messages.GetText("donation")
		keyboard = keyboards.BackToMainMenuKeyboard()
	case "contact":
		text = messages.GetText("contact")
		keyboard = keyboards.BackToMainMenuKeyboard()
	case "events":
		text = getEventsListText(ctx)
		keyboard = keyboards.BackToMainMenuKeyboard()

	default:
		text = messages.GetText("other_answer")
		keyboard = keyboards.BackToMainMenuKeyboard()
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

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		userID := update.Message.From.ID
		messageText := update.Message.Text

		conv := conversation.GetManager()
		state := conv.GetState(userID)

		if messageText == "/cancel" && state != "" {
			conv.ClearState(userID)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Операцію скасовано.",
			})
			return
		}

		if state != "" &&
			state != internalModels.StateAwaitingDeleteID &&
			state != internalModels.StateAwaitingDeleteConfirm &&
			state != internalModels.StateAwaitingBroadcastText &&
			state != internalModels.StateAwaitingBroadcastConfirm {
			HandleEventDialogMessage(ctx, b, update)
			return
		}

		if state == internalModels.StateAwaitingDeleteID && middleware.IsAdmin(userID) {
			DeleteEventHandler(ctx, b, update)
			return
		}

		if (state == internalModels.StateAwaitingBroadcastText ||
			state == internalModels.StateAwaitingBroadcastConfirm) &&
			middleware.IsAdmin(userID) {
			HandleBroadcastDialogMessage(ctx, b, update)
			return
		}

		var text string
		var keyboard *models.InlineKeyboardMarkup

		switch messageText {
		case "⛪ Про нас":
			text = messages.GetText("about_menu")
			keyboard = keyboards.AboutUsKeyboard()

		case "🙏 Служіння":
			text = messages.GetText("ministry_menu")
			keyboard = keyboards.MinistryKeyboard()

		case "📱 Соц. мережі":
			text = messages.GetText("social_media")
			keyboard = keyboards.BackToMainMenuKeyboard()

		case "📅 Події":
			text = getEventsListText(ctx)
			keyboard = keyboards.BackToMainMenuKeyboard()

		case "💳 Підтримати":
			text = messages.GetText("donation")
			keyboard = keyboards.BackToMainMenuKeyboard()

		case "📍 Наша адреса":
			text = messages.GetText("contact")
			keyboard = keyboards.BackToMainMenuKeyboard()

		case "🔕 Відписатися від розсилки":
			handleUnsubscribe(ctx, b, update)
			return

		case "🔔 Підписатися на розсилку":
			handleSubscribe(ctx, b, update)
			return

		default:
			text = messages.GetText("other_answer")

			isActive := true
			user, err := userRepo.GetByID(ctx, userID)
			if err == nil && user != nil {
				isActive = user.IsActive
			}

			replyKeyboard := keyboards.MainMenuReplyKeyboard(isActive)

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        text,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: replyKeyboard,
				LinkPreviewOptions: &models.LinkPreviewOptions{
					IsDisabled: bot.True(),
				},
			})

			return
		}

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

}

func getEventsListText(ctx context.Context) string {
	events, err := eventRepo.GetUpcoming(ctx)
	if err != nil {
		log.Printf("Error getting upcoming events: %v", err)
		return messages.GetText("no_events")
	}

	if len(events) == 0 {
		return messages.GetText("no_events")
	}

	text := "📅 <b>Найближчі події</b>\n\n"

	for i, event := range events {
		text += fmt.Sprintf(
			"<b>%d. %s</b>\n"+
				"📅 %s\n"+
				"📝 %s\n",
			i+1,
			event.Title,
			formatEventDate(event.Date),
			event.Description,
		)

		if event.Location != nil && *event.Location != "" {
			text += fmt.Sprintf("📍 %s\n", *event.Location)
		}

		if event.Category != nil && *event.Category != "" {
			text += fmt.Sprintf("🏷 %s\n", *event.Category)
		}

		if event.RegistrationURL != nil && *event.RegistrationURL != "" {
			text += fmt.Sprintf("🔗 <a href=\"%s\">Реєстрація</a>\n", *event.RegistrationURL)
		}

		text += "\n"
	}

	return text
}

func formatEventDate(t time.Time) string {
	if t.Hour() == 0 && t.Minute() == 0 {
		return t.Format("02.01.2006")
	}
	return t.Format("02.01.2006 15:04")
}

func handleUnsubscribe(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	err := userRepo.SetActive(ctx, userID, false)
	if err != nil {
		log.Printf("Error unsubscribing user: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Помилка відписки. Спробуйте пізніше.",
		})
		return
	}

	text := "🔕 <b>Ви відписалися від розсилки</b>\n\n" +
		"Ви більше не отримуватимете автоматичні повідомлення.\n" +
		"Але можете і далі користуватися ботом!\n\n" +
		"Щоб знову підписатися, натисніть кнопку нижче."

	keyboard := keyboards.MainMenuReplyKeyboard(false)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboard,
	})
}

func handleSubscribe(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	err := userRepo.SetActive(ctx, userID, true)
	if err != nil {
		log.Printf("Error subscribing user: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Помилка підписки. Спробуйте пізніше.",
		})
		return
	}

	text := "🔔 <b>Ви підписалися на розсилку!</b>\n\n" +
		"Тепер ви будете отримувати:\n" +
		"• Нагадування про богослужіння\n" +
		"• Інформацію про події\n" +
		"• Важливі оголошення\n\n" +
		"Ви завжди можете відписатися натиснувши кнопку нижче."

	keyboard := keyboards.MainMenuReplyKeyboard(true)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: keyboard,
	})
}
