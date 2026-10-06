package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/smford/cyber-hud/internal/config"
	"github.com/smford/cyber-hud/internal/theme"
)

type OverviewTabModel struct {
	config       *config.Config
	theme        *theme.Theme
	terminal     TerminalModel
	width        int
	height       int
	selectedCard int
}

func NewOverviewTabModel(cfg *config.Config, th *theme.Theme, w, h int) OverviewTabModel {
	termH := h - 12
	if termH < 10 {
		termH = 10
	}
	term := NewTerminalModel(cfg, th, w-2, termH)
	return OverviewTabModel{
		config:       cfg,
		theme:        th,
		terminal:     term,
		width:        w,
		height:       h,
		selectedCard: 0,
	}
}

func (m *OverviewTabModel) SetTheme(th *theme.Theme) {
	m.theme = th
	m.terminal.SetTheme(th)
}

func (m *OverviewTabModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	termH := h - 12
	if termH < 10 {
		termH = 10
	}
	m.terminal.SetSize(w-2, termH)
}

func (m OverviewTabModel) View() string {
	var sb strings.Builder

	// Netrunner Terminal directly at the top of the screen
	sb.WriteString(m.terminal.View() + "\n\n")

	// Recent Hardware Dispatches Header
	sectionTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("// RECENT HARDWARE DISPATCHES")
	line := lipgloss.NewStyle().Foreground(m.theme.Secondary).Render(strings.Repeat("─", 40))
	sb.WriteString(sectionTitle + " " + line + "\n")

	// Dispatch Cards (show top 2 from BuildLogs)
	var renderedCards []string
	cardW := (m.width - 8) / 2
	if cardW < 35 {
		cardW = 35
	}

	for i := 0; i < 2 && i < len(m.config.BuildLogs); i++ {
		log := m.config.BuildLogs[i]
		tag := m.theme.CardTag.Render(log.Category)
		meta := fmt.Sprintf("%s   %s", tag, lipgloss.NewStyle().Foreground(m.theme.TextMuted).Render(log.ID))
		cardTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render(log.Title)
		summary := lipgloss.NewStyle().Foreground(m.theme.TextMuted).Width(cardW - 4).Render(log.Summary)
		action := m.theme.ActionBtn.Render(fmt.Sprintf("[INSPECT SCHEMATICS (%s)]", log.ID))

		cardContent := lipgloss.JoinVertical(lipgloss.Left, meta, cardTitle, summary, action)
		cardStyle := m.theme.Card.Width(cardW)
		if i == m.selectedCard {
			cardStyle = m.theme.CardActive.Width(cardW)
		}
		renderedCards = append(renderedCards, cardStyle.Render(cardContent))
	}

	cardsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedCards...)
	sb.WriteString(cardsRow)

	return lipgloss.NewStyle().
		Padding(0, 1).
		Render(sb.String())
}
