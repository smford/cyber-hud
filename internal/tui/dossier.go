package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smford/cyber-hud/internal/config"
	"github.com/smford/cyber-hud/internal/theme"
)

type DossierModel struct {
	log      *config.BuildLog
	theme    *theme.Theme
	viewport viewport.Model
	width    int
	height   int
	active   bool
}

func NewDossierModel(th *theme.Theme, w, h int) DossierModel {
	vp := viewport.New(w-8, h-8)
	return DossierModel{
		theme:    th,
		viewport: vp,
		width:    w,
		height:   h,
		active:   false,
	}
}

func (m *DossierModel) SetTheme(th *theme.Theme) {
	m.theme = th
}

func (m *DossierModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w - 8
	m.viewport.Height = h - 10
	if m.log != nil {
		m.renderContent()
	}
}

func (m *DossierModel) Open(log *config.BuildLog) {
	m.log = log
	m.active = true
	m.renderContent()
}

func (m *DossierModel) Close() {
	m.active = false
	m.log = nil
}

func (m *DossierModel) renderContent() {
	if m.log == nil {
		return
	}
	var sb strings.Builder

	// Title & Meta
	title := lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true).Render(fmt.Sprintf("// TECHNICAL SCHEMATIC DOSSIER: %s", strings.ToUpper(m.log.Title)))
	meta := fmt.Sprintf("[%s]  |  DATE: %s  |  STATUS: %s", m.log.ID, m.log.Date, m.log.Status)
	sb.WriteString(title + "\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(m.theme.TextMuted).Render(meta) + "\n\n")

	// Summary
	sb.WriteString(lipgloss.NewStyle().Foreground(m.theme.TextMain).Render(m.log.Summary) + "\n\n")

	// Schematic Diagram
	if m.log.Schematic != "" {
		sb.WriteString(lipgloss.NewStyle().Foreground(m.theme.Secondary).Bold(true).Render("▶ HARDWARE TOPOLOGY & SCHEMATIC FLOW:\n"))
		box := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(m.theme.Secondary).
			Padding(0, 1).
			Render(m.log.Schematic)
		sb.WriteString(box + "\n\n")
	}

	// Bill of Materials (BOM) Table
	if len(m.log.BOM) > 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true).Render("▶ BILL OF MATERIALS (BOM):\n"))
		sb.WriteString(fmt.Sprintf("  %-36s %-6s %-28s\n", "COMPONENT", "QTY", "FUNCTION / ROLE"))
		sb.WriteString("  " + strings.Repeat("─", 72) + "\n")
		for _, item := range m.log.BOM {
			sb.WriteString(fmt.Sprintf("  %-36s %-6s %-28s\n", item.Item, item.Qty, item.Role))
		}
		sb.WriteString("\n")
	}

	// Pinout Table
	if len(m.log.Pinout) > 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(m.theme.Secondary).Bold(true).Render("▶ GPIO & PINOUT SPECIFICATION:\n"))
		sb.WriteString(fmt.Sprintf("  %-16s %-22s %-32s\n", "GPIO / PIN", "TARGET DEVICE", "FUNCTIONAL NOTE"))
		sb.WriteString("  " + strings.Repeat("─", 72) + "\n")
		for _, p := range m.log.Pinout {
			sb.WriteString(fmt.Sprintf("  %-16s %-22s %-32s\n", p.GPIO, p.Device, p.Note))
		}
		sb.WriteString("\n")
	}

	// Engineering Notes
	if m.log.Notes != "" {
		sb.WriteString(lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true).Render("▶ FIELD TELEMETRY & ENGINEERING NOTES:\n"))
		sb.WriteString(lipgloss.NewStyle().Foreground(m.theme.TextMuted).Render(m.log.Notes) + "\n")
	}

	m.viewport.SetContent(sb.String())
	m.viewport.GotoTop()
}

func (m DossierModel) Update(msg tea.Msg) (DossierModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			m.Close()
			return m, nil
		case "up", "k":
			m.viewport.LineUp(2)
			return m, nil
		case "down", "j":
			m.viewport.LineDown(2)
			return m, nil
		case "pgup":
			m.viewport.HalfViewUp()
			return m, nil
		case "pgdown":
			m.viewport.HalfViewDown()
			return m, nil
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			if msg.Button == tea.MouseButtonWheelUp {
				m.viewport.LineUp(3)
				return m, nil
			}
			if msg.Button == tea.MouseButtonWheelDown {
				m.viewport.LineDown(3)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m DossierModel) View() string {
	if !m.active || m.log == nil {
		return ""
	}

	closeBtn := lipgloss.NewStyle().
		Background(m.theme.Danger).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1).
		Render("[ESC] CLOSE_DOSSIER")

	modalHeader := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true).Render(fmt.Sprintf("DOSSIER // %s", m.log.ID)),
		strings.Repeat(" ", 8),
		closeBtn,
	)

	body := m.viewport.View()

	content := lipgloss.JoinVertical(lipgloss.Left, modalHeader, "\n", body)

	modalWidth := m.width - 6
	if modalWidth < 50 {
		modalWidth = 50
	}
	modalHeight := m.height - 4
	if modalHeight < 15 {
		modalHeight = 15
	}

	return lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		Background(m.theme.BgPanel).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(m.theme.Primary).
		Padding(1, 2).
		Render(content)
}
