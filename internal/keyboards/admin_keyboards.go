package keyboards

import "github.com/go-telegram/bot/models"

func AdminPanelKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "➕ Додати подію", CallbackData: "admin_add_event"},
			},
			{
				{Text: "📋 Список подій", CallbackData: "admin_list_events"},
			},
			{
				{Text: "🏠 Головне меню", CallbackData: "back_to_start"},
			},
		},
	}
}

func AdminEventsListKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "➕ Додати подію", CallbackData: "admin_add_event"},
			},
			{
				{Text: "🗑️ Видалити подію", CallbackData: "admin_delete_event"},
			},
			{
				{Text: "◀️ Назад", CallbackData: "admin_panel"},
				{Text: "🏠 Головне меню", CallbackData: "back_to_start"},
			},
		},
	}
}

func BackToAdminPanelKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "◀️ До адмін-панелі", CallbackData: "admin_panel"},
			},
		},
	}
}

func DeleteConfirmKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Так, видалити", CallbackData: "admin_confirm_delete"},
				{Text: "❌ Скасувати", CallbackData: "admin_cancel_delete"},
			},
		},
	}
}
