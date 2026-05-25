package main

import (
	"fmt"
	"os"
	"os/exec"
)

func playAudio() error {
	c, err := loadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if c.AudioPath == "" {
		return nil
	}
	if _, err := os.Stat(c.AudioPath); err != nil {
		return fmt.Errorf("audio file not accessible: %w", err)
	}
	afplay, err := exec.LookPath("afplay")
	if err != nil {
		return fmt.Errorf("afplay not found in PATH (devtag requires macOS)")
	}
	cmd := exec.Command(afplay, c.AudioPath)
	return cmd.Run()
}
