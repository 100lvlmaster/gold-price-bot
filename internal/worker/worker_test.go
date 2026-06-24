package worker

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type mockTelegramBot struct {
	mu           sync.Mutex
	sentMessages []tgbotapi.Chattable
}

func (m *mockTelegramBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sentMessages = append(m.sentMessages, c)
	return tgbotapi.Message{}, nil
}

func (m *mockTelegramBot) getSentMessages() []tgbotapi.Chattable {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Return a copy to avoid race conditions
	copied := make([]tgbotapi.Chattable, len(m.sentMessages))
	copy(copied, m.sentMessages)
	return copied
}
