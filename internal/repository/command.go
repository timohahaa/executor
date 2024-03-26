package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/timohahaa/executor/internal/entity"

	"github.com/timohahaa/postgres"
)

type commandRepository struct {
	db *postgres.Postgres
}

func NewCommandRepository(pg *postgres.Postgres) *commandRepository {
	return &commandRepository{
		db: pg,
	}
}

func (r *commandRepository) CreateCommand(ctx context.Context, commandText string) (entity.Command, error) {
	sql, args, _ := r.db.Builder.
		Insert("commands").
		Columns("command_text").
		Values(commandText).
		Suffix("RETURNING command_id").
		ToSql()

	var createdId uint64
	err := r.db.ConnPool.QueryRow(ctx, sql, args...).Scan(&createdId)
	if err != nil {
		return entity.Command{}, err
	}

	return entity.Command{Id: createdId, Text: commandText}, nil
}

func (r *commandRepository) ListCommands(ctx context.Context, limit, offset uint64) ([]entity.Command, error) {
	sql, args, _ := r.db.Builder.
		Select("command_id", "command_text").
		From("commands").
		Limit(limit).
		Offset(offset).
		ToSql()

	var commands []entity.Command
	rows, err := r.db.ConnPool.Query(ctx, sql, args...)
	if err != nil {
		return []entity.Command(nil), err
	}

	for rows.Next() {
		var command entity.Command
		err = rows.Scan(&command)
		commands = append(commands, command)
	}

	return commands, err
}

func (r *commandRepository) GetCommandById(ctx context.Context, commandId uint64) (entity.Command, error) {
	sql, args, _ := r.db.Builder.
		Select("command_id", "command_text").
		From("commands").
		Where("command_id = ?", commandId).
		ToSql()

	var command entity.Command
	err := r.db.ConnPool.QueryRow(ctx, sql, args...).Scan(&command)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Command{}, ErrCommandNotFound
	} else if err != nil {
		return entity.Command{}, err
	}

	return command, nil
}
