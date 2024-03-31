package service

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/timohahaa/executor/internal/mocks/repomocks"
)

// TODO:
// ListCommands(ctx context.Context, limit, offset uint64) ([]entity.Command, error)
// StopCommand(ctx context.Context, commandId uint64) error

func TestCreateCommand(t *testing.T) {
	testCases := []struct {
		name        string
		cmdText     string
		wantCmdText string
		output      string
		wantOutput  string
		err         error
		wantErr     bool
	}{
		{
			name:        "command properely created",
			cmdText:     "ls -la",
			wantCmdText: "ls -la",
			output:      "",
			wantOutput:  "",
			err:         nil,
			wantErr:     false,
		},
		{
			name:        "error properely returned",
			cmdText:     "rm -r /",
			wantCmdText: "rm -r /",
			output:      "",
			wantOutput:  "",
			err:         errors.New("some error"),
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// создаем моки и все необходимое
			mockRepo := repomocks.NewCommandRepoMock()
			mockRepo.SetError(tc.err)
			log := logrus.New()
			log.SetOutput(io.Discard)

			// создаем сервис
			s := NewCommandService(mockRepo, "/bin/sh", log)

			// run
			cmd, err := s.CreateCommand(context.TODO(), tc.cmdText)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.wantCmdText, cmd.Text)
			assert.Equal(t, tc.wantOutput, cmd.LastOutput)
			assert.NoError(t, err)
		})
	}
}

func TestDeletCommandById(t *testing.T) {
	testCases := []struct {
		name       string
		cmdId      uint64
		cmdIdToSet uint64
		err        error
		wantErr    bool
	}{
		{
			name:       "deleting existing command",
			cmdId:      8,
			cmdIdToSet: 8,
			err:        nil,
			wantErr:    false,
		},
		{
			name:       "deleting non-existing command",
			cmdId:      28,
			cmdIdToSet: 8,
			err:        nil,
			wantErr:    false,
		},
		{
			name:    "error properly returned",
			cmdId:   1337,
			err:     errors.New("some error"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// создаем моки и все необходимое
			mockRepo := repomocks.NewCommandRepoMock()
			mockRepo.SetError(tc.err)
			mockRepo.SetCommand(tc.cmdIdToSet, "", "")
			log := logrus.New()
			log.SetOutput(io.Discard)

			// создаем сервис
			s := NewCommandService(mockRepo, "/bin/sh", log)

			// run
			err := s.DeleteCommandById(context.TODO(), tc.cmdId)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.False(t, mockRepo.HasCommand(tc.cmdId))
		})
	}
}

func TestGetCommandById(t *testing.T) {
	testCases := []struct {
		name        string
		cmdId       uint64
		cmdIdToSet  uint64
		cmdText     string
		wantCmdText string
		output      string
		wantOutput  string
		err         error
		wantErr     bool
	}{
		{
			name:        "get existing command",
			cmdId:       8,
			cmdIdToSet:  8,
			cmdText:     "ls -la",
			wantCmdText: "ls -la",
			output:      "",
			wantOutput:  "",
			err:         nil,
			wantErr:     false,
		},
		{
			name:        "get non-existing command",
			cmdId:       28,
			cmdIdToSet:  8,
			cmdText:     "ls -la",
			wantCmdText: "ls -la",
			output:      "",
			wantOutput:  "",
			err:         nil,
			wantErr:     true,
		},
		{
			name:        "error properly returned",
			cmdId:       1337,
			cmdIdToSet:  1337,
			cmdText:     "ls -la",
			wantCmdText: "ls -la",
			output:      "",
			wantOutput:  "",
			err:         errors.New("some error"),
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// создаем моки и все необходимое
			mockRepo := repomocks.NewCommandRepoMock()
			mockRepo.SetError(tc.err)
			mockRepo.SetCommand(tc.cmdIdToSet, tc.cmdText, tc.output)
			log := logrus.New()
			log.SetOutput(io.Discard)

			// создаем сервис
			s := NewCommandService(mockRepo, "/bin/sh", log)

			// Run
			cmd, err := s.GetCommandById(context.TODO(), tc.cmdId)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCmdText, cmd.Text)
			assert.Equal(t, tc.wantOutput, cmd.LastOutput)
		})
	}
}

func TestRunCommand(t *testing.T) {
	testCases := []struct {
		name     string
		cmdId    uint64
		pid      int
		pidToSet int
		err      error
		wantErr  bool
	}{
		{
			name:     "run already running command",
			cmdId:    8,
			pid:      228,
			pidToSet: 228,
			err:      nil,
			wantErr:  true,
		},
		{
			name:     "error properly returned",
			cmdId:    8,
			pid:      28,
			pidToSet: 228,
			err:      errors.New("some error"),
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// создаем моки и все необходимое
			mockRepo := repomocks.NewCommandRepoMock()
			mockRepo.SetError(tc.err)
			mockRepo.SetPID(tc.cmdId, tc.pidToSet)
			log := logrus.New()
			log.SetOutput(io.Discard)

			// создаем сервис
			s := NewCommandService(mockRepo, "/bin/sh", log)

			// Run
			err := s.RunCommand(context.TODO(), tc.cmdId)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
