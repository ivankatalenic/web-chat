package services

import (
	"github.com/ivankatalenic/web-chat/impl/logger"
	"github.com/ivankatalenic/web-chat/impl/repository"
	"github.com/ivankatalenic/web-chat/interfaces"
)

// GetMessageRepository getter
func GetMessageRepository() interfaces.MessageRepository {
	return repository.NewInMemory()
}

func GetLogger() interfaces.Logger {
	return logger.NewConsoleLogger()
}
