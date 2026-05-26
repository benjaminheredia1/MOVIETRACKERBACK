package application

import (
	"MovieTrackerBack/internal/domain"
	"errors"
)

type ChatService struct {
	repo domain.ChatRepository
}

func NewChatService(repo domain.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) SaveMessage(msg domain.ChatMessage) (string, error) {
	if msg.Content == "" {
		return "", errors.New("message content cannot be empty")
	}
	if msg.SessionID == "" {
		return "", errors.New("session id cannot be empty")
	}

	reply, err := s.repo.GenerateMessage(msg)
	if err != nil {
		return "", err
	}

	return reply, nil
}
