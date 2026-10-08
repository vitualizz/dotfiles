package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type InstallLogger struct {
	file    *os.File
	mu      sync.Mutex
	logPath string
}

func NewInstallLogger() (*InstallLogger, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		stateDir = filepath.Join(home, ".local", "state")
	}
	logDir := filepath.Join(stateDir, "vitualizz-dotfiles")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, "install.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(f, "\n=== Session started: %s ===\n", time.Now().Format("2006-01-02 15:04:05"))

	return &InstallLogger{file: f, logPath: logPath}, nil
}

func (l *InstallLogger) LogPath() string {
	return l.logPath
}

func (l *InstallLogger) LogCommand(toolName, command string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "[%s] CMD [%s]: %s\n", timestamp(), toolName, command)
	_ = l.file.Sync()
}

func (l *InstallLogger) LogSuccess(toolName string, duration time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "[%s] OK   [%s] (%s)\n", timestamp(), toolName, duration)
	_ = l.file.Sync()
}

func (l *InstallLogger) LogError(toolName, command string, err error, output string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "[%s] ERR  [%s]\n", timestamp(), toolName)
	fmt.Fprintf(l.file, "       Command: %s\n", command)
	fmt.Fprintf(l.file, "       Error:   %v\n", err)
	if output != "" {
		out := truncate(output, 2000)
		fmt.Fprintf(l.file, "       Output:\n%s\n", indent(out, "       "))
	}
	fmt.Fprintln(l.file)
	_ = l.file.Sync()
}

func (l *InstallLogger) LogInfo(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "[%s] INFO: %s\n", timestamp(), msg)
	_ = l.file.Sync()
}

func (l *InstallLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func timestamp() string {
	return time.Now().Format("15:04:05")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... [truncated]"
}

func indent(s, prefix string) string {
	return prefix + strings.ReplaceAll(s, "\n", "\n"+prefix)
}
