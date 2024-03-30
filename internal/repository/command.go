package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/timohahaa/executor/internal/entity"

	"github.com/timohahaa/postgres"
)

type commandRepository struct {
	db  *postgres.Postgres
	rdb *redis.Client
	// как долно хранить pid исполняемой команды (в секундах)
	pidTTL int64
	log    *logrus.Logger
}

func NewCommandRepository(pg *postgres.Postgres, redis *redis.Client, pidTTL int64, logger *logrus.Logger) *commandRepository {
	return &commandRepository{
		db:     pg,
		rdb:    redis,
		pidTTL: pidTTL,
		log:    logger,
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
		r.log.Errorf("commandRepository.CreateCommand -> QueryRow: %v", err)
		return entity.Command{}, err
	}

	return entity.Command{Id: createdId, Text: commandText}, nil
}

func (r *commandRepository) ListCommands(ctx context.Context, limit, offset uint64) ([]entity.Command, error) {
	sql, args, _ := r.db.Builder.
		Select("command_id", "command_text", "last_output").
		From("commands").
		Limit(limit).
		Offset(offset).
		ToSql()

	var commands []entity.Command
	rows, err := r.db.ConnPool.Query(ctx, sql, args...)
	if err != nil {
		r.log.Errorf("commandRepository.ListCommands -> Query: %v", err)
		return []entity.Command(nil), err
	}

	for rows.Next() {
		var command entity.Command
		err = rows.Scan(&command.Id, &command.Text, &command.LastOutput)
		commands = append(commands, command)
	}

	return commands, err
}

func (r *commandRepository) GetCommandById(ctx context.Context, commandId uint64) (entity.Command, error) {
	sql, args, _ := r.db.Builder.
		Select("command_id", "command_text", "last_output").
		From("commands").
		Where("command_id = ?", commandId).
		ToSql()

	var command entity.Command
	err := r.db.ConnPool.QueryRow(ctx, sql, args...).Scan(&command.Id, &command.Text, &command.LastOutput)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Command{}, ErrCommandNotFound
	} else if err != nil {
		r.log.Errorf("commandRepository.GetCommandById -> QueryRow: %v", err)
		return entity.Command{}, err
	}

	return command, nil
}

func (r *commandRepository) SaveCommandOutput(ctx context.Context, commandId uint64, line string) error {
	// нет поддержки set с сложными присвоениями у squirrel
	sql := "UPDATE commands SET last_output = last_output || $1 WHERE command_id = $2"
	args := []interface{}{line, commandId}

	_, err := r.db.ConnPool.Exec(ctx, sql, args...)
	if err != nil {
		r.log.Errorf("commandRepository.SaveCommandOutput -> Exec: %v", err)
		return err
	}

	return nil
}
func (r *commandRepository) ClearCommandOutput(ctx context.Context, commandId uint64) error {
	sql, args, _ := r.db.Builder.
		Update("commands").
		Set("last_output", "").
		Where("command_id = ?", commandId).
		ToSql()

	_, err := r.db.ConnPool.Exec(ctx, sql, args...)
	if err != nil {
		r.log.Errorf("commandRepository.ClearCommandOutput -> Exec: %v", err)
		return err
	}

	return nil
}

func (r *commandRepository) SetCommandPID(ctx context.Context, commandId uint64, pid int) error {
	key := fmt.Sprintf("command:%d", commandId)
	err := r.rdb.Set(ctx, key, pid, time.Second*time.Duration(r.pidTTL)).Err()
	if err != nil {
		r.log.Errorf("commandRepository.SetCommandPID -> Set: %v", err)
		return err
	}

	return nil
}

func (r *commandRepository) GetCommandPID(ctx context.Context, commandId uint64) (int, error) {
	key := fmt.Sprintf("command:%d", commandId)
	val, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, ErrCommandNotRunning
	} else if err != nil {
		r.log.Errorf("commandRepository.GetCommandPID -> Get: %v", err)
		return 0, err
	}

	// поскольку всегда сохраняем int, не возникнет ошибки парсинга
	pid, _ := strconv.Atoi(val)

	return pid, nil
}

func (r *commandRepository) DeleteCommandPID(ctx context.Context, commandId uint64) error {
	key := fmt.Sprintf("command:%d", commandId)
	_, err := r.rdb.Del(ctx, key).Result()
	if err != nil {
		r.log.Errorf("commandRepository.DeleteCommandPID -> Del: %v", err)
		return err
	}

	return nil
}
