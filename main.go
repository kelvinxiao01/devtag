package main

import (
	"fmt"
	"os"
)

const usage = `devtag — play an audio clip whenever you git push

Usage:
  devtag set <path>      Set the audio file played on git push
  devtag get             Print the currently configured audio path
  devtag play            Play the configured audio (used by the git hook)
  devtag install         Install the global git pre-push hook
  devtag install --force   Overwrite an existing core.hooksPath
  devtag uninstall       Remove the hook and unset core.hooksPath
  devtag --help          Show this message
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Print(usage)
		return
	}

	switch args[0] {
	case "set":
		if len(args) < 2 {
			die("devtag set requires a path argument")
		}
		runSet(args[1])
	case "get":
		runGet()
	case "play":
		if err := playAudio(); err != nil {
			fmt.Fprintln(os.Stderr, "devtag play:", err)
		}
	case "install":
		force := false
		for _, a := range args[1:] {
			if a == "--force" || a == "-f" {
				force = true
			}
		}
		if err := installHooks(force); err != nil {
			die(err.Error())
		}
	case "uninstall":
		if err := uninstallHooks(); err != nil {
			die(err.Error())
		}
	default:
		die(fmt.Sprintf("unknown command %q\n\n%s", args[0], usage))
	}
}

func runSet(p string) {
	abs, err := absExpand(p)
	if err != nil {
		die(err.Error())
	}
	info, err := os.Stat(abs)
	if err != nil {
		die(fmt.Sprintf("cannot access %s: %v", abs, err))
	}
	if info.IsDir() {
		die(fmt.Sprintf("%s is a directory, not a file", abs))
	}
	if err := saveConfig(config{AudioPath: abs}); err != nil {
		die(err.Error())
	}
	fmt.Printf("Audio path set to %s\n", abs)
}

func runGet() {
	c, err := loadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "No audio configured. Run: devtag set <path>")
			os.Exit(1)
		}
		die(err.Error())
	}
	if c.AudioPath == "" {
		fmt.Fprintln(os.Stderr, "No audio configured. Run: devtag set <path>")
		os.Exit(1)
	}
	fmt.Println(c.AudioPath)
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
