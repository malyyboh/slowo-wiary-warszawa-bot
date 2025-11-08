package keyboards

import (
	"github.com/go-telegram/bot/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/messages"
)

func AboutUsKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: messages.AboutUsButtons["about_church"], CallbackData: "about_church"},
				{Text: messages.AboutUsButtons["about_church_mission"], CallbackData: "about_church_mission"},
			},
			{
				{Text: messages.AboutUsButtons["about_church_doctrine"], CallbackData: "about_church_doctrine"},
				{Text: messages.AboutUsButtons["about_church_pastors"], CallbackData: "about_church_pastors"},
			},
			{
				{Text: messages.AboutUsButtons["about_church_history"], CallbackData: "about_church_history"},
				{Text: messages.NavigationButtons["back_to_start"], CallbackData: "back_to_start"},
			},
		},
	}
}

func MinistryKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: messages.MinistryButtons["sunday_ministry"], CallbackData: "sunday_ministry"},
				{Text: messages.MinistryButtons["home_ministry"], CallbackData: "home_ministry"},
			},
			{
				{Text: messages.MinistryButtons["prayer_ministry"], CallbackData: "prayer_ministry"},
				{Text: messages.MinistryButtons["youth_ministry"], CallbackData: "youth_ministry"},
			},
			{
				{Text: messages.MinistryButtons["teenagers_ministry"], CallbackData: "teenagers_ministry"},
				{Text: messages.MinistryButtons["kindergarten_ministry"], CallbackData: "kindergarten_ministry"},
			},
			{
				{Text: messages.NavigationButtons["back_to_start"], CallbackData: "back_to_start"},
			},
		},
	}
}

func BackToMainMenuKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: messages.NavigationButtons["back_to_start"], CallbackData: "back_to_start"},
			},
		},
	}
}

func BackToAboutKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: messages.NavigationButtons["back"], CallbackData: "about_us"},
				{Text: messages.NavigationButtons["back_to_start"], CallbackData: "back_to_start"},
			},
		},
	}
}

func BackToMinistryKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: messages.NavigationButtons["back"], CallbackData: "ministry"},
				{Text: messages.NavigationButtons["back_to_start"], CallbackData: "back_to_start"},
			},
		},
	}
}

func MainMenuReplyKeyboard(isActive bool) *models.ReplyKeyboardMarkup {
	var buttons [][]models.KeyboardButton

	buttons = append(buttons, []models.KeyboardButton{
		{Text: messages.MainMenuButtons["about_us"]},
		{Text: messages.MainMenuButtons["ministry"]},
	})

	buttons = append(buttons, []models.KeyboardButton{
		{Text: messages.MainMenuButtons["social_media"]},
		{Text: messages.MainMenuButtons["events"]},
	})

	buttons = append(buttons, []models.KeyboardButton{
		{Text: messages.MainMenuButtons["donation"]},
		{Text: messages.MainMenuButtons["contact"]},
	})

	if isActive {
		buttons = append(buttons, []models.KeyboardButton{
			{Text: "🔕 Відписатися від розсилки"},
		})
	} else {
		buttons = append(buttons, []models.KeyboardButton{
			{Text: "🔔 Підписатися на розсилку"},
		})
	}

	return &models.ReplyKeyboardMarkup{
		Keyboard:       buttons,
		ResizeKeyboard: true,
	}
}
