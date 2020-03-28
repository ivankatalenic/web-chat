package services

import (
	"github.com/ivankatalenic/web-chat/internal/impl/logger"
	"github.com/ivankatalenic/web-chat/internal/impl/repository"
	"github.com/ivankatalenic/web-chat/internal/interfaces"
)

// GetMessageRepository getter
func GetMessageRepository() interfaces.MessageRepository {
	return repository.NewInMemory()
}

func GetLogger() interfaces.Logger {
	return logger.NewConsoleLogger()
}
