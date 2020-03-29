package services

import (
	"github.com/ivankatalenic/web-chat/internal/impl/logger"
	"github.com/ivankatalenic/web-chat/internal/impl/repository"
	"github.com/ivankatalenic/web-chat/internal/interfaces"
)

// GetMessageRepository getter
func GetMessageRepository() interfaces.MessageRepository {
	const repositorySize = 1024
	return repository.NewInMemory(repositorySize)
}

func GetLogger() interfaces.Logger {
	return logger.NewConsoleLogger()
}
