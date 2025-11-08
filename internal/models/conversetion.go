package models

type ConversationState struct {
	UserID        int64
	State         string
	EventData     *Event
	BroadcastText string
}

const (
	StateIdle                     = ""
	StateAwaitingTitle            = "awaiting_title"
	StateAwaitingDate             = "awaiting_date"
	StateAwaitingDesc             = "awaiting_description"
	StateAwaitingLocation         = "awaiting_location"
	StateAwaitingCategory         = "awaiting_category"
	StateAwaitingRegURL           = "awaiting_registration_url"
	StateAwaitingConfirm          = "awaiting_confirmation"
	StateAwaitingDeleteID         = "awaiting_delete_id"
	StateAwaitingDeleteConfirm    = "awaiting_delete_confirm"
	StateAwaitingBroadcastText    = "awaiting_broadcast_text"
	StateAwaitingBroadcastConfirm = "awaiting_broadcast_confirm"
)
