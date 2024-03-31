package repomocks

import (
	"context"
	"math/rand"

	"github.com/timohahaa/executor/internal/entity"
	"github.com/timohahaa/executor/internal/repository"
)

type CommandRepoMock struct {
	cmdMap   map[uint64]*entity.Command
	cmdSlice []entity.Command
	pidMap   map[uint64]int
	err      error
}

var _ repository.CommandRepository = &CommandRepoMock{}

func NewCommandRepoMock() *CommandRepoMock {
	return &CommandRepoMock{
		cmdMap:   make(map[uint64]*entity.Command),
		pidMap:   make(map[uint64]int),
		cmdSlice: make([]entity.Command, 0),
		err:      nil,
	}
}

func (m *CommandRepoMock) makeCommand(id uint64, text, output string) entity.Command {
	return entity.Command{
		Id:         id,
		Text:       text,
		LastOutput: output,
	}
}

func (m *CommandRepoMock) SetCommand(id uint64, text, output string) {
	cmd := m.makeCommand(id, text, output)
	m.cmdMap[cmd.Id] = &cmd
	m.cmdSlice = append(m.cmdSlice, cmd)
}

func (m *CommandRepoMock) SetError(err error) {
	m.err = err
}

func (m *CommandRepoMock) SetPID(id uint64, pid int) {
	m.pidMap[id] = rand.Int()
}

func (m *CommandRepoMock) HasCommand(id uint64) bool {
	_, ok := m.cmdMap[id]
	return ok
}

func (m *CommandRepoMock) CreateCommand(_ context.Context, commandText string) (entity.Command, error) {
	cmd := m.makeCommand(rand.Uint64(), commandText, "")
	m.cmdMap[cmd.Id] = &cmd
	m.cmdSlice = append(m.cmdSlice, cmd)
	return cmd, m.err
}
func (m *CommandRepoMock) DeleteCommandById(_ context.Context, commandId uint64) error {
	delete(m.cmdMap, commandId)
	for i := range m.cmdSlice {
		if m.cmdSlice[i].Id == commandId {
			m.cmdSlice = append(m.cmdSlice[:i], m.cmdSlice[i+1:]...)
			break
		}
	}
	return m.err
}
func (m *CommandRepoMock) ListCommands(_ context.Context, limit, offset uint64) ([]entity.Command, error) {
	if offset > uint64(len(m.cmdSlice)) {
		return nil, m.err
	}
	cmds := m.cmdSlice[offset:]
	if limit >= uint64(len(cmds)) {
		return cmds, m.err
	}
	return cmds[:limit], m.err
}
func (m *CommandRepoMock) GetCommandById(_ context.Context, commandId uint64) (entity.Command, error) {
	cmd, ok := m.cmdMap[commandId]
	if !ok {
		return entity.Command{}, repository.ErrCommandNotFound
	}
	return *cmd, m.err
}
func (m *CommandRepoMock) SaveCommandOutput(_ context.Context, commandId uint64, line string) error {
	cmd, ok := m.cmdMap[commandId]
	if !ok {
		return m.err
	}
	cmd.LastOutput += line
	return m.err
}
func (m *CommandRepoMock) ClearCommandOutput(_ context.Context, commandId uint64) error {
	cmd, ok := m.cmdMap[commandId]
	if !ok {
		return m.err
	}
	cmd.LastOutput = ""
	return m.err
}
func (m *CommandRepoMock) SetCommandPID(_ context.Context, commandId uint64, pid int) error {
	m.pidMap[commandId] = pid
	return m.err
}
func (m *CommandRepoMock) GetCommandPID(_ context.Context, commandId uint64) (int, error) {
	pid, ok := m.pidMap[commandId]
	if !ok {
		return 0, m.err
	}
	return pid, m.err
}
func (m *CommandRepoMock) DeleteCommandPID(_ context.Context, commandId uint64) error {
	delete(m.pidMap, commandId)
	return m.err
}
