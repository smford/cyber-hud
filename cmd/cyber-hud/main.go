package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smford/cyber-hud/internal/config"
	"github.com/smford/cyber-hud/internal/tui"
)

var (
	version = "1.0.0"
)

func main() {
	configPath := flag.String("c", "", "Path to YAML or JSON config file (default: embedded/local)")
	fullscreen := flag.Bool("F", false, "Request full screen terminal mode (macOS)")
	skipBoot := flag.Bool("skip-boot", false, "Skip boot animation directly into HUD")
	dumpConfig := flag.String("dump-config", "", "Export default configuration to specified YAML/JSON file and exit")
	showVersion := flag.Bool("v", false, "Display cyber-hud version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("cyber-hud v%s (Go %s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// If user requested dumping default configuration
	if *dumpConfig != "" {
		cfg := config.DefaultConfig()
		if err := config.Save(*dumpConfig, &cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Default configuration successfully exported to: %s\n", *dumpConfig)
		os.Exit(0)
	}

	// Check if local config.yaml exists if not explicitly passed
	targetConfig := *configPath
	if targetConfig == "" {
		if _, err := os.Stat("config.yaml"); err == nil {
			targetConfig = "config.yaml"
		} else if _, err := os.Stat("config.json"); err == nil {
			targetConfig = "config.json"
		}
	}

	// Load configuration
	cfg, err := config.Load(targetConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load %s (%v). Using embedded defaults.\n", targetConfig, err)
		defaultCfg := config.DefaultConfig()
		cfg = &defaultCfg
	}

	if *skipBoot {
		cfg.Settings.SkipBoot = true
	}

	// Optional: Request native fullscreen window on macOS
	if *fullscreen && runtime.GOOS == "darwin" {
		script := `tell application "System Events" to tell (first process whose frontmost is true) to keystroke "f" using {control down, command down}`
		_ = exec.Command("osascript", "-e", script).Run()
	}

	// Initialize Bubble Tea TUI
	appModel := tui.NewAppModel(cfg)

	// Options: AltScreen takes over the full terminal display, mouse tracking enabled
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
	}

	if cfg.Settings.MouseEnabled {
		opts = append(opts, tea.WithMouseCellMotion())
	}

	p := tea.NewProgram(appModel, opts...)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running cyber-hud: %v\n", err)
		os.Exit(1)
	}
}
