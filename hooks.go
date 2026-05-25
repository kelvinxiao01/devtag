package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const prePushHook = `#!/bin/sh
# devtag: play configured audio on git push (non-blocking)
(devtag play >/dev/null 2>&1 &) >/dev/null 2>&1
exit 0
`

func hooksDir() (string, error) {
	dir, err := devtagDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "hooks"), nil
}

func gitGlobalHooksPath() (string, error) {
	out, err := exec.Command("git", "config", "--global", "--get", "core.hooksPath").Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func installHooks(force bool) error {
	hd, err := hooksDir()
	if err != nil {
		return err
	}
	expandedHd, err := absExpand(hd)
	if err != nil {
		return err
	}

	existing, err := gitGlobalHooksPath()
	if err != nil {
		return fmt.Errorf("read git core.hooksPath: %w", err)
	}
	existingAbs, _ := absExpand(existing)
	if existing != "" && existingAbs != expandedHd && !force {
		return fmt.Errorf(
			"git core.hooksPath is already set to %q; re-run with `devtag install --force` to overwrite",
			existing,
		)
	}

	if err := os.MkdirAll(hd, 0o755); err != nil {
		return err
	}
	hookPath := filepath.Join(hd, "pre-push")
	if err := os.WriteFile(hookPath, []byte(prePushHook), 0o755); err != nil {
		return err
	}

	if err := exec.Command("git", "config", "--global", "core.hooksPath", hd).Run(); err != nil {
		return fmt.Errorf("set git core.hooksPath: %w", err)
	}

	fmt.Printf("Installed pre-push hook at %s\n", hookPath)
	fmt.Printf("Set git --global core.hooksPath = %s\n", hd)
	fmt.Println("Note: while devtag is installed, git ignores per-repo .git/hooks/ for hooks.")
	return nil
}

func uninstallHooks() error {
	hd, err := hooksDir()
	if err != nil {
		return err
	}
	expandedHd, err := absExpand(hd)
	if err != nil {
		return err
	}

	hookPath := filepath.Join(hd, "pre-push")
	if err := os.Remove(hookPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	existing, err := gitGlobalHooksPath()
	if err != nil {
		return fmt.Errorf("read git core.hooksPath: %w", err)
	}
	existingAbs, _ := absExpand(existing)
	if existing != "" && existingAbs == expandedHd {
		if err := exec.Command("git", "config", "--global", "--unset", "core.hooksPath").Run(); err != nil {
			return fmt.Errorf("unset git core.hooksPath: %w", err)
		}
		fmt.Println("Unset git --global core.hooksPath")
	} else if existing != "" {
		fmt.Printf("Left git core.hooksPath alone (points to %q, not devtag's hooks dir)\n", existing)
	}
	fmt.Printf("Removed %s\n", hookPath)
	return nil
}

func absExpand(p string) (string, error) {
	if p == "" {
		return "", nil
	}
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(home, strings.TrimPrefix(p, "~"))
	}
	return filepath.Abs(p)
}
