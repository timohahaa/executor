package service

import "errors"

var (
	ErrCommandNotFound       = errors.New("command not found")
	ErrCommandNotRunning     = errors.New("command not currently running")
	ErrCommandAlreadyRunning = errors.New("command already running")
)
