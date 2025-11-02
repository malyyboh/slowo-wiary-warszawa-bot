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
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/repository"
)

var eventRepo = repository.NewEventRepository()

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := messages.GetText("/start")
	keyboard := keyboards.MainMenuKeyboard()

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
	keyboard := keyboards.MainMenuKeyboard()

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

func CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	data := callback.Data

	var text string
	var keyboard *models.InlineKeyboardMarkup

	switch data {
	case "back_to_start":
		text = messages.GetText("/start")
		keyboard = keyboards.MainMenuKeyboard()

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
		text = getEventsListText()
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

		conv := conversation.GetManager()
		state := conv.GetState(userID)

		if state != "" {
			HandleEventDialogMessage(ctx, b, update)
			return
		}

		text := messages.GetText("other_answer")
		keyboard := keyboards.MainMenuKeyboard()

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

func getEventsListText() string {
	events, err := eventRepo.GetUpcoming()
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
