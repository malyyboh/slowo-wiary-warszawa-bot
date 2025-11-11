package keyboards

import "github.com/go-telegram/bot/models"

func AdminPanelKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "➕ Додати подію", CallbackData: "admin_add_event"},
				{Text: "📋 Список подій", CallbackData: "admin_list_events"},
			},
			{
				{Text: "📊 Користувачі", CallbackData: "admin_users"},
				{Text: "💾 Експорт БД", CallbackData: "admin_export_db"},
			},
			{
				{Text: "📢 Розсилка", CallbackData: "admin_broadcast"},
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

func AdminUsersKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "📋 Список користувачів", CallbackData: "admin_list_users"},
			},
			{
				{Text: "◀️ Назад", CallbackData: "admin_panel"},
				{Text: "🏠 Головне меню", CallbackData: "back_to_start"},
			},
		},
	}
}

func AdminUsersListKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "◀️ Назад", CallbackData: "admin_users"},
				{Text: "🏠 Головне меню", CallbackData: "back_to_start"},
			},
		},
	}
}

func AdminBroadcastKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "📤 Надіслати зараз", CallbackData: "admin_broadcast_now"},
			},
			{
				{Text: "◀️ Назад", CallbackData: "admin_panel"},
				{Text: "🏠 Головне меню", CallbackData: "back_to_start"},
			},
		},
	}
}

func BroadcastConfirmKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Так, надіслати", CallbackData: "admin_confirm_broadcast"},
				{Text: "❌ Скасувати", CallbackData: "admin_cancel_broadcast"},
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
