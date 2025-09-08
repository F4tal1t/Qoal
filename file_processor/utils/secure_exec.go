package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type SecureCommandExecutor struct {
	timeout time.Duration
}

func NewSecureCommandExecutor(timeout time.Duration) *SecureCommandExecutor {
	return &SecureCommandExecutor{timeout: timeout}
}

func (e *SecureCommandExecutor) ExecuteCommand(name string, args []string) error {
	// Validate command name (whitelist approach)
	if !e.isAllowedCommand(name) {
		return fmt.Errorf("command not allowed: %s", name)
	}

	// Sanitize arguments
	sanitizedArgs := e.sanitizeArgs(args)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// Execute command
	cmd := exec.CommandContext(ctx, name, sanitizedArgs...)

	return cmd.Run()
}

func (e *SecureCommandExecutor) isAllowedCommand(name string) bool {
	allowedCommands := map[string]bool{
		"ffmpeg":      true,
		"convert":     true, // ImageMagick
		"libreoffice": true,
		"7z":          true,
		"unzip":       true,
		"clamscan":    true,
	}

	// Extract base command name
	baseName := strings.TrimSuffix(name, ".exe")
	if idx := strings.LastIndex(baseName, "/"); idx != -1 {
		baseName = baseName[idx+1:]
	}
	if idx := strings.LastIndex(baseName, "\\"); idx != -1 {
		baseName = baseName[idx+1:]
	}

	return allowedCommands[baseName]
}

func (e *SecureCommandExecutor) sanitizeArgs(args []string) []string {
	sanitized := make([]string, len(args))

	for i, arg := range args {
		// Remove potentially dangerous characters
		arg = strings.ReplaceAll(arg, ";", "")
		arg = strings.ReplaceAll(arg, "&", "")
		arg = strings.ReplaceAll(arg, "|", "")
		arg = strings.ReplaceAll(arg, "`", "")
		arg = strings.ReplaceAll(arg, "$", "")

		sanitized[i] = arg
	}

	return sanitized
}
