package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func FormatCode(ctx context.Context, filePath string, buildTag string) error {
	goplsPath, err := exec.LookPath("gopls")
	if err != nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("os.UserHomeDir() failed: %w", err)
		}
		goplsPath = filepath.Join(homeDir, "go", "bin", "gopls")
		if len(goplsPath) <= 0 {
			return fmt.Errorf("executable file \"gopls\" not found")
		}
	}
	cmd := exec.CommandContext(ctx, goplsPath, "format", "-w", filePath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s failed: %w", cmd.String(), err)
	}
	cmd = exec.CommandContext(ctx, goplsPath, "imports", "-w", filePath)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOFLAGS=\"-tags=%s\"", buildTag))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s failed: %w", cmd.String(), err)
	}
	return nil
}
