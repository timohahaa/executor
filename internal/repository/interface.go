package repository

import (
	"context"

	"github.com/timohahaa/executor/internal/entity"
)

type CommandRepository interface {
	CreateCommand(ctx context.Context, commandText string) (entity.Command, error)
	ListCommands(ctx context.Context, limit, offset uint64) ([]entity.Command, error)
	GetCommandById(ctx context.Context, commandId uint64) (entity.Command, error)
}

var _ CommandRepository = &commandRepository{}
