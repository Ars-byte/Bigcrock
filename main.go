package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"bigcrock/src"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version", "version":
			fmt.Printf("bigcrock %s\n", src.Version)
			return
		case "-h", "--help", "help":
			fmt.Println("Usage: bigcrock [directory]")
			fmt.Println("\nA fast, Ranger-like TUI directory tree & disk usage ranking analyzer.")
			fmt.Println("\nOptions:")
			fmt.Println("  -v, --version  Print version and exit")
			fmt.Println("  -h, --help     Print help and exit")
			return
		}
	}

	target := "."
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	target = strings.Trim(target, "\"'`")
	if target == "~" || strings.HasPrefix(target, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			if target == "~" {
				target = home
			} else {
				target = filepath.Join(home, target[2:])
			}
		}
	}

	info, err := os.Stat(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot access %q: %v\n", target, err)
		os.Exit(1)
	}

	if !info.IsDir() {
		target = filepath.Dir(target)
	}

	m := src.NewModel(target)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running bigcrock: %v\n", err)
		os.Exit(1)
	}
}
