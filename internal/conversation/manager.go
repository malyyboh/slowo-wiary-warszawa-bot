package conversation

import (
	"sync"

	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
)

type Manager struct {
	conversations map[int64]*models.ConversationState
	mu            sync.RWMutex
}

var manager *Manager

func InitManager() {
	manager = &Manager{
		conversations: make(map[int64]*models.ConversationState),
	}
}

func GetManager() *Manager {
	return manager
}

func (m *Manager) SetState(userID int64, state string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conv, exists := m.conversations[userID]; exists {
		conv.State = state
	} else {
		m.conversations[userID] = &models.ConversationState{
			UserID:    userID,
			State:     state,
			EventData: &models.Event{},
		}
	}
}

func (m *Manager) GetState(userID int64) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if conv, exists := m.conversations[userID]; exists {
		return conv.State
	}
	return models.StateIdle
}

func (m *Manager) GetConversation(userID int64) *models.ConversationState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if conv, exists := m.conversations[userID]; exists {
		return conv
	}
	return nil
}

func (m *Manager) ClearState(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.conversations, userID)
}
