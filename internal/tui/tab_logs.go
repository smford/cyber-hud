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

type LogsTabModel struct {
	config       *config.Config
	theme        *theme.Theme
	width        int
	height       int
	categories   []string
	activeCatIdx int
	selectedIdx  int
	viewport     viewport.Model
	triggerModal string
}

func NewLogsTabModel(cfg *config.Config, th *theme.Theme, w, h int) LogsTabModel {
	cats := []string{"ALL LOGS", "Hardware", "Optics", "Fabrication"}
	vp := viewport.New(w-4, h-8)
	return LogsTabModel{
		config:       cfg,
		theme:        th,
		width:        w,
		height:       h,
		categories:   cats,
		activeCatIdx: 0,
		selectedIdx:  0,
		viewport:     vp,
	}
}

func (m *LogsTabModel) SetTheme(th *theme.Theme) {
	m.theme = th
}

func (m *LogsTabModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w - 4
	m.viewport.Height = h - 8
	if m.viewport.Height < 5 {
		m.viewport.Height = 5
	}
}

func (m LogsTabModel) filteredLogs() []config.BuildLog {
	cat := m.categories[m.activeCatIdx]
	if cat == "ALL LOGS" {
		return m.config.BuildLogs
	}
	var filtered []config.BuildLog
	for _, l := range m.config.BuildLogs {
		if strings.EqualFold(l.Category, cat) {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

func (m LogsTabModel) Update(msg tea.Msg) (LogsTabModel, tea.Cmd) {
	logs := m.filteredLogs()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil
		case "down", "j":
			if m.selectedIdx < len(logs)-1 {
				m.selectedIdx++
			}
			return m, nil
		case "left", "h":
			if m.activeCatIdx > 0 {
				m.activeCatIdx--
				m.selectedIdx = 0
			}
			return m, nil
		case "right", "l":
			if m.activeCatIdx < len(m.categories)-1 {
				m.activeCatIdx++
				m.selectedIdx = 0
			}
			return m, nil
		case "enter", "space":
			if len(logs) > 0 && m.selectedIdx < len(logs) {
				m.triggerModal = logs[m.selectedIdx].ID
			}
			return m, nil
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			if msg.Button == tea.MouseButtonWheelUp {
				m.viewport.LineUp(2)
				return m, nil
			}
			if msg.Button == tea.MouseButtonWheelDown {
				m.viewport.LineDown(2)
				return m, nil
			}
		}
	}
	return m, nil
}

func (m LogsTabModel) View() string {
	var sb strings.Builder

	// Section Head
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("// HARDWARE ENGINEERING BUILD LOGS")
	line := lipgloss.NewStyle().Foreground(m.theme.Secondary).Render(strings.Repeat("─", 40))
	sb.WriteString(title + " " + line + "\n\n")

	// Filter Bar
	var renderedFilters []string
	for i, cat := range m.categories {
		if i == m.activeCatIdx {
			renderedFilters = append(renderedFilters, m.theme.ActionBtnHover.Render(cat))
		} else {
			renderedFilters = append(renderedFilters, m.theme.ActionBtn.Render(cat))
		}
	}
	filterBar := lipgloss.JoinHorizontal(lipgloss.Center, renderedFilters...)
	sb.WriteString(filterBar + "\n\n")

	// Logs Grid
	logs := m.filteredLogs()
	cardW := (m.width - 8) / 2
	if cardW < 36 {
		cardW = 36
	}

	var renderedCards []string
	for i, log := range logs {
		tag := m.theme.CardTag.Render(log.Category)
		meta := fmt.Sprintf("%s  %s  |  DATE: %s  |  %s", tag, log.ID, log.Date, log.Status)
		cTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render(log.Title)
		summary := lipgloss.NewStyle().Foreground(m.theme.TextMuted).Width(cardW - 4).Render(log.Summary)

		action := m.theme.ActionBtn.Render(fmt.Sprintf("[INSPECT SCHEMATICS & BOM (%s)]", log.ID))
		if i == m.selectedIdx {
			action = m.theme.ActionBtnHover.Render(fmt.Sprintf("[INSPECT SCHEMATICS & BOM (%s)] (ENTER)", log.ID))
		}

		cardContent := lipgloss.JoinVertical(lipgloss.Left, meta, cTitle, summary, action)
		cardStyle := m.theme.Card.Width(cardW)
		if i == m.selectedIdx {
			cardStyle = m.theme.CardActive.Width(cardW)
		}
		renderedCards = append(renderedCards, cardStyle.Render(cardContent))
	}

	// Pair cards into 2-column rows
	var cardRows []string
	for i := 0; i < len(renderedCards); i += 2 {
		if i+1 < len(renderedCards) {
			row := lipgloss.JoinHorizontal(lipgloss.Top, renderedCards[i], "  ", renderedCards[i+1])
			cardRows = append(cardRows, row)
		} else {
			cardRows = append(cardRows, renderedCards[i])
		}
	}

	gridContent := strings.Join(cardRows, "\n\n")
	m.viewport.SetContent(gridContent)
	sb.WriteString(m.viewport.View())

	return lipgloss.NewStyle().Padding(0, 1).Render(sb.String())
}
