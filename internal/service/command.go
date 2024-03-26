package service

import (
	"context"
	"errors"

	"github.com/timohahaa/executor/internal/entity"
	"github.com/timohahaa/executor/internal/repository"
)

type commandService struct {
	commandRepo repository.CommandRepository
}

func NewCommandService(commandRepo repository.CommandRepository) *commandService {
	return &commandService{
		commandRepo: commandRepo,
	}
}

func (s *commandService) CreateCommand(ctx context.Context, commandText string) (entity.Command, error) {
	command, err := s.commandRepo.CreateCommand(ctx, commandText)
	if err != nil {
		return entity.Command{}, err
	}

	return command, nil
}

func (s *commandService) ListCommands(ctx context.Context, limit, offset uint64) ([]entity.Command, error) {
	commands, err := s.commandRepo.ListCommands(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (s *commandService) GetCommandById(ctx context.Context, commandId uint64) (entity.Command, error) {
	command, err := s.commandRepo.GetCommandById(ctx, commandId)
	if errors.Is(err, repository.ErrCommandNotFound) {
		return entity.Command{}, ErrCommandNotFound
	} else if err != nil {
		return entity.Command{}, err
	}

	return command, nil
}
