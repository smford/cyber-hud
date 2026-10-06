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

type FirmwareTabModel struct {
	config       *config.Config
	theme        *theme.Theme
	width        int
	height       int
	categories   []string
	activeCatIdx int
	selectedIdx  int
	viewport     viewport.Model
}

func NewFirmwareTabModel(cfg *config.Config, th *theme.Theme, w, h int) FirmwareTabModel {
	cats := []string{"ALL REPOSITORIES", "Embedded", "RaspberryPi", "Networking", "3DPrinting", "Terminal", "Web", "Fabrication"}
	vp := viewport.New(w-4, h-8)
	return FirmwareTabModel{
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

func (m *FirmwareTabModel) SetTheme(th *theme.Theme) {
	m.theme = th
}

func (m *FirmwareTabModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w - 4
	m.viewport.Height = h - 8
	if m.viewport.Height < 5 {
		m.viewport.Height = 5
	}
}

func (m FirmwareTabModel) filteredRepos() []config.Repository {
	cat := m.categories[m.activeCatIdx]
	if cat == "ALL REPOSITORIES" {
		return m.config.Repositories
	}
	var filtered []config.Repository
	for _, r := range m.config.Repositories {
		if strings.EqualFold(r.Category, cat) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func (m FirmwareTabModel) Update(msg tea.Msg) (FirmwareTabModel, tea.Cmd) {
	repos := m.filteredRepos()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil
		case "down", "j":
			if m.selectedIdx < len(repos)-1 {
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
			if len(repos) > 0 && m.selectedIdx < len(repos) {
				repo := repos[m.selectedIdx]
				targetURL := repo.URL
				if repo.Demo != "" {
					targetURL = repo.Demo
				}
				openBrowser(targetURL)
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

func (m FirmwareTabModel) View() string {
	var sb strings.Builder

	// Section Head
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).
		Render(fmt.Sprintf("// GITHUB REPOSITORY REGISTRY (%d NODES)", len(m.config.Repositories)))
	line := lipgloss.NewStyle().Foreground(m.theme.Secondary).Render(strings.Repeat("─", 36))
	sb.WriteString(title + " " + line + "\n\n")

	// Filter Bar
	var renderedFilters []string
	for i, cat := range m.categories {
		label := cat
		if i == 0 {
			label = fmt.Sprintf("ALL (%d)", len(m.config.Repositories))
		}
		if i == m.activeCatIdx {
			renderedFilters = append(renderedFilters, m.theme.ActionBtnHover.Render(label))
		} else {
			renderedFilters = append(renderedFilters, m.theme.ActionBtn.Render(label))
		}
	}
	filterBar := lipgloss.JoinHorizontal(lipgloss.Center, renderedFilters...)
	sb.WriteString(filterBar + "\n\n")

	// Repositories List
	repos := m.filteredRepos()
	cardW := (m.width - 8) / 2
	if cardW < 36 {
		cardW = 36
	}

	var renderedCards []string
	for i, r := range repos {
		tag := m.theme.CardTag.Render(r.Category)
		starsBadge := lipgloss.NewStyle().Foreground(m.theme.Primary).Render(fmt.Sprintf("★ %d", r.Stars))
		forksBadge := lipgloss.NewStyle().Foreground(m.theme.Secondary).Render(fmt.Sprintf("⑂ %d", r.Forks))
		meta := fmt.Sprintf("%s   %s   %s", tag, starsBadge, forksBadge)

		rTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render(r.Name)
		desc := lipgloss.NewStyle().Foreground(m.theme.TextMuted).Width(cardW - 4).Render(r.Desc)

		var pills []string
		for _, t := range r.Tech {
			pills = append(pills, m.theme.TechPill.Render(t))
		}
		techLine := lipgloss.JoinHorizontal(lipgloss.Center, pills...)

		actionLabel := "[OPEN REPO / DEMO]"
		if i == m.selectedIdx {
			actionLabel = "[OPEN REPO / DEMO] (ENTER)"
		}
		action := m.theme.ActionBtn.Render(actionLabel)
		if i == m.selectedIdx {
			action = m.theme.ActionBtnHover.Render(actionLabel)
		}

		cardContent := lipgloss.JoinVertical(lipgloss.Left, meta, rTitle, desc, techLine, action)
		cardStyle := m.theme.Card.Width(cardW)
		if i == m.selectedIdx {
			cardStyle = m.theme.CardActive.Width(cardW)
		}
		renderedCards = append(renderedCards, cardStyle.Render(cardContent))
	}

	// Pair cards into 2-column grid
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
