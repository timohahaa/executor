package service

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/timohahaa/executor/internal/entity"
	"github.com/timohahaa/executor/internal/repository"
)

type commandService struct {
	commandRepo repository.CommandRepository
	// пусть до стандартной оболочки, например /bin/sh или /bin/bash
	defaultShellPath string
}

func NewCommandService(commandRepo repository.CommandRepository, defaultShellPath string) *commandService {
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
	// либо команда не найдена, либо ошибка на уровне подключения к БД
	if errors.Is(err, repository.ErrCommandNotFound) {
		return entity.Command{}, ErrCommandNotFound
	} else if err != nil {
		return entity.Command{}, err
	}

	return command, nil
}

func (s *commandService) processOutput(ctx context.Context, commandId uint64, line string) error {
	err := s.commandRepo.SaveCommandOutput(ctx, commandId, line)
	if err != nil {
		return err
	}

	return nil
}

func (s *commandService) RunCommand(ctx context.Context, commandId uint64) error {
	// запущена ли уже команда?
	_, err := s.commandRepo.GetCommandPID(ctx, commandId)
	// команда уже запущена
	if err == nil {
		return ErrCommandAlreadyRunning
	}
	// не можем проверить, запущена ли команда, так как возникла ошибка
	if err != nil && !errors.Is(err, repository.ErrCommandNotRunning) {
		return err
	}

	command, err := s.commandRepo.GetCommandById(ctx, commandId)
	// либо команда не найдена, либо ошибка на уровне подключения к БД
	if errors.Is(err, repository.ErrCommandNotFound) {
		return ErrCommandNotFound
	} else if err != nil {
		return err
	}

	cmd := exec.Command(s.defaultShellPath, "-c", command.Text)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// поскольку ставим false, сигналы к программе будут дублироваться на дочерние процессы
		// НО зато будут права на то, чтобы эти дочерние процессы останавливать
		Setpgid: false,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)

	// стартуем исполнение команды
	err = cmd.Start()
	if err != nil {
		return err
	}

	// сохраняем  pid команды
	s.commandRepo.SetCommandPID(ctx, commandId, cmd.Process.Pid)
	// при завершении команды или возникновении ошибки "забываем" её pid
	defer func() {
		s.commandRepo.DeleteCommandPID(ctx, commandId)
	}()

	// сканируем и сохраняем вывод команды
	for scanner.Scan() {
		_ = s.processOutput(ctx, commandId, scanner.Text())
	}

	if scanner.Err() != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return scanner.Err()
	}

	return cmd.Wait()
}

func (s *commandService) StopCommand(ctx context.Context, commandId uint64) error {
	// проверяем наличие команды
	_, err := s.commandRepo.GetCommandById(ctx, commandId)
	// либо команда не найдена, либо ошибка на уровне подключения к БД
	if errors.Is(err, repository.ErrCommandNotFound) {
		return ErrCommandNotFound
	} else if err != nil {
		return err
	}

	// получаем pid команды
	pid, err := s.commandRepo.GetCommandPID(ctx, commandId)
	if errors.Is(err, repository.ErrCommandNotRunning) {
		return ErrCommandNotRunning
	} else if err != nil {
		return err
	}

	// после завершения должны удалить pid команды
	defer func() {
		s.commandRepo.DeleteCommandPID(ctx, commandId)
	}()

	// нет ошибки на unix системах
	process, _ := os.FindProcess(pid)
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		// процесс уже был завершен
		return nil
	}

	// попробуем завершить по хорошему :)
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		// по хорошему не хочет, значит будет по плохому >:(
		err = process.Signal(syscall.SIGKILL)
		return err
	}

	// проверяем что процесс завершился
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			// процесс завершился
			return nil
		}
	}

	// по хорошему он снова не захотел >:(
	err = process.Signal(syscall.SIGKILL)
	return err
}
