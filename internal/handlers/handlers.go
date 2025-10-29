package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/keyboards"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/messages"
)

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
		text = messages.GetText("no_events")
		keyboard = keyboards.BackToMainMenuKeyboard()

	default:
		text = messages.GetText("other_answer")
		keyboard = keyboards.MainMenuKeyboard()
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
