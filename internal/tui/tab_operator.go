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

type OperatorTabModel struct {
	config   *config.Config
	theme    *theme.Theme
	width    int
	height   int
	viewport viewport.Model
}

func NewOperatorTabModel(cfg *config.Config, th *theme.Theme, w, h int) OperatorTabModel {
	vp := viewport.New(w-4, h-8)
	return OperatorTabModel{
		config:   cfg,
		theme:    th,
		width:    w,
		height:   h,
		viewport: vp,
	}
}

func (m *OperatorTabModel) SetTheme(th *theme.Theme) {
	m.theme = th
}

func (m *OperatorTabModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w - 4
	m.viewport.Height = h - 8
	if m.viewport.Height < 5 {
		m.viewport.Height = 5
	}
}

func (m OperatorTabModel) Update(msg tea.Msg) (OperatorTabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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

func (m OperatorTabModel) View() string {
	var sb strings.Builder

	// Section Head: Skills
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("// OPERATOR PROFILE & CYBERWARE SPEC")
	line := lipgloss.NewStyle().Foreground(m.theme.Secondary).Render(strings.Repeat("─", 40))
	sb.WriteString(title + " " + line + "\n\n")

	// Skills Grid
	skillBarW := 36
	if m.width > 80 {
		skillBarW = 48
	}

	var renderedSkills []string
	for _, skill := range m.config.Skills {
		filled := (skill.Percentage * skillBarW) / 100
		empty := skillBarW - filled
		bar := lipgloss.NewStyle().Foreground(m.theme.Primary).Render(strings.Repeat("█", filled)) +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#1a2533")).Render(strings.Repeat("░", empty))

		header := fmt.Sprintf("%-34s %s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render(skill.Name),
			lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true).Render(fmt.Sprintf("%3d%%", skill.Percentage)))

		skillBox := lipgloss.JoinVertical(lipgloss.Left, header, bar)
		card := m.theme.Card.Width(skillBarW + 6).Render(skillBox)
		renderedSkills = append(renderedSkills, card)
	}

	// Pair skills into 2-column layout
	var skillRows []string
	for i := 0; i < len(renderedSkills); i += 2 {
		if i+1 < len(renderedSkills) {
			row := lipgloss.JoinHorizontal(lipgloss.Top, renderedSkills[i], "  ", renderedSkills[i+1])
			skillRows = append(skillRows, row)
		} else {
			skillRows = append(skillRows, renderedSkills[i])
		}
	}
	sb.WriteString(strings.Join(skillRows, "\n\n") + "\n\n")

	// Section Head: Benches
	benchTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("// WORKSHOP HARDWARE INVENTORY")
	sb.WriteString(benchTitle + " " + line + "\n\n")

	benchCardW := (m.width - 12) / 3
	if benchCardW < 30 {
		benchCardW = 30
	}

	var renderedBenches []string
	for _, b := range m.config.Benches {
		bHeader := fmt.Sprintf("[%s] %s", b.ID, b.Title)
		headerStyled := lipgloss.NewStyle().Foreground(m.theme.Secondary).Bold(true).Render(bHeader)

		var items []string
		for _, it := range b.Items {
			bullet := lipgloss.NewStyle().Foreground(m.theme.Primary).Render("•")
			items = append(items, fmt.Sprintf("%s %s", bullet, lipgloss.NewStyle().Foreground(m.theme.TextMuted).Render(it)))
		}
		itemsContent := strings.Join(items, "\n")
		cardContent := lipgloss.JoinVertical(lipgloss.Left, headerStyled, "\n", itemsContent)
		card := m.theme.Card.Width(benchCardW).Render(cardContent)
		renderedBenches = append(renderedBenches, card)
	}

	benchesRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedBenches...)
	sb.WriteString(benchesRow + "\n\n")

	// Footer Note
	footerNote := lipgloss.NewStyle().Foreground(m.theme.TextMuted).
		Render("EXECUTED VIA NEURAL DEBIAN KERNEL // MAINTAINED BY @smford // NIGHT CITY PROTOCOL")
	sb.WriteString(footerNote)

	m.viewport.SetContent(sb.String())
	return lipgloss.NewStyle().Padding(0, 1).Render(m.viewport.View())
}
