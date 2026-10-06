package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/smford/cyber-hud/internal/theme"
)

type TabID int

const (
	TabOverview TabID = iota
	TabBuildLogs
	TabFirmware
	TabOperator
)

type TabItem struct {
	ID    TabID
	Key   string
	Title string
	Badge int
}

type TabsModel struct {
	activeTab TabID
	theme     *theme.Theme
	width     int
	items     []TabItem
}

func NewTabsModel(th *theme.Theme, logCount, repoCount int) TabsModel {
	return TabsModel{
		activeTab: TabOverview,
		theme:     th,
		items: []TabItem{
			{ID: TabOverview, Key: "1", Title: "OVERVIEW_HUD", Badge: 0},
			{ID: TabBuildLogs, Key: "2", Title: "BUILD_LOGS", Badge: logCount},
			{ID: TabFirmware, Key: "3", Title: "FIRMWARE_DECK", Badge: repoCount},
			{ID: TabOperator, Key: "4", Title: "OPERATOR_STACK", Badge: 0},
		},
	}
}

func (m *TabsModel) SetActive(tab TabID) {
	m.activeTab = tab
}

func (m *TabsModel) Active() TabID {
	return m.activeTab
}

func (m *TabsModel) Next() {
	m.activeTab = (m.activeTab + 1) % TabID(len(m.items))
}

func (m *TabsModel) Prev() {
	if m.activeTab == 0 {
		m.activeTab = TabID(len(m.items) - 1)
	} else {
		m.activeTab--
	}
}

func (m *TabsModel) SetTheme(th *theme.Theme) {
	m.theme = th
}

func (m *TabsModel) SetWidth(w int) {
	m.width = w
}

func (m TabsModel) View() string {
	var renderedTabs []string

	for _, item := range m.items {
		var tabStr string
		label := fmt.Sprintf("[%s] %s", item.Key, item.Title)
		if item.Badge > 0 {
			badgeStr := m.theme.TabBadge.Render(fmt.Sprintf("%d", item.Badge))
			label = fmt.Sprintf("%s %s", label, badgeStr)
		}

		if item.ID == m.activeTab {
			tabStr = m.theme.TabActive.Render(label)
		} else {
			tabStr = m.theme.TabInactive.Render(label)
		}
		renderedTabs = append(renderedTabs, tabStr)
	}

	leftBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	themeBtn := lipgloss.NewStyle().
		Background(lipgloss.Color("#101924")).
		Foreground(m.theme.Primary).
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.theme.Secondary).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf("STYLE: [ %s ] (F1)", m.theme.Name))

	gap := m.width - lipgloss.Width(leftBar) - lipgloss.Width(themeBtn) - 4
	var fullBar string
	if gap > 2 {
		fullBar = leftBar + strings.Repeat(" ", gap) + themeBtn
	} else {
		fullBar = leftBar
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Background(m.theme.BgBase).
		Padding(0, 1).
		Render(fullBar)
}

// CheckMouseClick returns the TabID if (x, y) clicked on a tab, or -1 if not.
func (m TabsModel) CheckMouseClick(x, y int) (TabID, bool) {
	// Tabs start at y=0 now that the top bar is removed
	if y != 0 && y != 1 {
		return 0, false
	}
	currentX := 1
	for _, item := range m.items {
		label := fmt.Sprintf("[%s] %s", item.Key, item.Title)
		if item.Badge > 0 {
			label = fmt.Sprintf("%s %d", label, item.Badge)
		}
		tabWidth := len(label) + 4 // padding & margins
		if x >= currentX && x < currentX+tabWidth {
			return item.ID, true
		}
		currentX += tabWidth + 1
	}
	return 0, false
}

// CheckThemeClick returns true if mouse clicked the theme button on the right of the tab bar.
func (m TabsModel) CheckThemeClick(x, y int) bool {
	if y != 0 && y != 1 {
		return false
	}
	themeBtnLen := len(m.theme.Name) + 18
	return x >= m.width-themeBtnLen-4
}
