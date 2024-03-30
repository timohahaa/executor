package repository

import "errors"

var (
	ErrCommandNotFound   = errors.New("command not found")
	ErrCommandNotRunning = errors.New("command is not currently running")
)
