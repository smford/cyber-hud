package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smford/cyber-hud/internal/config"
	"github.com/smford/cyber-hud/internal/theme"
)

func TestTerminalCommands(t *testing.T) {
	cfg := config.DefaultConfig()
	th := theme.Cyberpunk2077()
	term := NewTerminalModel(&cfg, th, 100, 30)

	// Test help command
	term.execute("help")
	lastLine := term.outputLines[len(term.outputLines)-1]
	if !strings.Contains(lastLine, "INTERACTIVE TELEMETRY") {
		t.Errorf("expected help output, got %s", lastLine)
	}

	// Test projects command
	term.execute("projects")
	found := false
	for _, l := range term.outputLines {
		if strings.Contains(l, "GITHUB REPOSITORY INDEX") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected projects list output")
	}

	// Test cat VFS file
	term.execute("cat notes.md")
	foundNotes := false
	for _, l := range term.outputLines {
		if strings.Contains(l, "Night City Hardware Bench") {
			foundNotes = true
			break
		}
	}
	if !foundNotes {
		t.Errorf("expected notes.md content")
	}

	// Test CIDR calculation
	term.execute("cidr 10.0.0.0/24")
	foundCidr := false
	for _, l := range term.outputLines {
		if strings.Contains(l, "Total Addresses: 256") {
			foundCidr = true
			break
		}
	}
	if !foundCidr {
		t.Errorf("expected CIDR report")
	}

	// Test custom command
	term.execute("esp32")
	foundESP := false
	for _, l := range term.outputLines {
		if strings.Contains(l, "ESP32 ASYNC WEBSERVER") {
			foundESP = true
			break
		}
	}
	if !foundESP {
		t.Errorf("expected esp32 output")
	}
}

func TestTTY1OnlineVisible(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Settings.SkipBoot = true
	app := NewAppModel(&cfg)
	m, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 35})
	view := m.View()
	lines := strings.Split(view, "\n")
	if len(lines) > 35 {
		t.Fatalf("expected <= 35 lines, got %d", len(lines))
	}
	found := false
	for i := 0; i < 4 && i < len(lines); i++ {
		if strings.Contains(lines[i], "TTY1_ONLINE") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected TTY1_ONLINE in the top lines of view")
	}
}
