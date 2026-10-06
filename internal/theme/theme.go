package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines the visual palette and styles for the HUD.
type Theme struct {
	ID        string
	Name      string
	Primary   lipgloss.Color // Yellow / Ice Blue
	Secondary lipgloss.Color // Cyan / Violet
	Danger    lipgloss.Color // Red / Hot Pink
	Success   lipgloss.Color // Green
	Purple    lipgloss.Color // Purple / Deep violet
	BgBase    lipgloss.Color // Dark base background
	BgPanel   lipgloss.Color // Panel background
	TextMain  lipgloss.Color // Foreground text
	TextMuted lipgloss.Color // Muted gray/blue

	// Pre-computed Lip Gloss styles
	TopBar         lipgloss.Style
	StatusActive   lipgloss.Style
	TabActive      lipgloss.Style
	TabInactive    lipgloss.Style
	TabBadge       lipgloss.Style
	PanelBox       lipgloss.Style
	Card           lipgloss.Style
	CardActive     lipgloss.Style
	CardTag        lipgloss.Style
	TechPill       lipgloss.Style
	ActionBtn      lipgloss.Style
	ActionBtnHover lipgloss.Style
	ChipBtn        lipgloss.Style
	ChipBtnActive  lipgloss.Style
	SectionTitle   lipgloss.Style
	SectionLine    lipgloss.Style
	PromptStyle    lipgloss.Style
	CommandEcho    lipgloss.Style
	ErrorOutput    lipgloss.Style
	ModalBox       lipgloss.Style
}

// Cyberpunk2077 returns the Night City yellow/cyan aesthetic.
func Cyberpunk2077() *Theme {
	t := &Theme{
		ID:        "cyberpunk",
		Name:      "CYBERPUNK 2077",
		Primary:   lipgloss.Color("#FCEE0A"), // Night City Yellow
		Secondary: lipgloss.Color("#00F0FF"), // Neon Cyan
		Danger:    lipgloss.Color("#FF003C"), // Arasaka Red
		Success:   lipgloss.Color("#00FF66"), // Matrix Green
		Purple:    lipgloss.Color("#9D4EDD"),
		BgBase:    lipgloss.Color("#06080E"),
		BgPanel:   lipgloss.Color("#0D121C"),
		TextMain:  lipgloss.Color("#DCF2FD"),
		TextMuted: lipgloss.Color("#6B8299"),
	}
	t.initStyles()
	return t
}

// AlteredCarbon returns the Psychasec ice-blue/needlecast violet aesthetic.
func AlteredCarbon() *Theme {
	t := &Theme{
		ID:        "altered-carbon",
		Name:      "ALTERED CARBON",
		Primary:   lipgloss.Color("#00E5FF"), // Ice Blue
		Secondary: lipgloss.Color("#BD00FF"), // Needlecast Violet
		Danger:    lipgloss.Color("#FF2A6D"),
		Success:   lipgloss.Color("#05FFA1"),
		Purple:    lipgloss.Color("#7928CA"),
		BgBase:    lipgloss.Color("#030509"),
		BgPanel:   lipgloss.Color("#080C16"),
		TextMain:  lipgloss.Color("#E2F1F8"),
		TextMuted: lipgloss.Color("#577087"),
	}
	t.initStyles()
	return t
}

func (t *Theme) initStyles() {
	t.TopBar = lipgloss.NewStyle().
		Background(t.BgPanel).
		Foreground(t.TextMain).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(t.Secondary)

	t.StatusActive = lipgloss.NewStyle().
		Foreground(t.Success).
		Bold(true)

	t.TabActive = lipgloss.NewStyle().
		Background(t.Secondary).
		Foreground(lipgloss.Color("#000000")).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	t.TabInactive = lipgloss.NewStyle().
		Background(t.BgPanel).
		Foreground(t.TextMuted).
		Padding(0, 1).
		MarginRight(1).
		Border(lipgloss.NormalBorder(), false, false, false, false)

	t.TabBadge = lipgloss.NewStyle().
		Background(t.Danger).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1).
		MarginLeft(1)

	t.PanelBox = lipgloss.NewStyle().
		Background(t.BgPanel).
		Foreground(t.TextMain).
		Padding(1, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Secondary)

	t.Card = lipgloss.NewStyle().
		Background(t.BgPanel).
		Foreground(t.TextMain).
		Padding(1, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#1B2A38"))

	t.CardActive = lipgloss.NewStyle().
		Background(t.BgPanel).
		Foreground(t.TextMain).
		Padding(1, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary)

	t.CardTag = lipgloss.NewStyle().
		Background(t.Primary).
		Foreground(lipgloss.Color("#000000")).
		Bold(true).
		Padding(0, 1)

	t.TechPill = lipgloss.NewStyle().
		Foreground(t.Secondary).
		Background(lipgloss.Color("#132230")).
		Padding(0, 1).
		MarginRight(1)

	t.ActionBtn = lipgloss.NewStyle().
		Background(t.BgPanel).
		Foreground(t.Secondary).
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.Secondary).
		Padding(0, 1)

	t.ActionBtnHover = lipgloss.NewStyle().
		Background(t.Secondary).
		Foreground(lipgloss.Color("#000000")).
		Bold(true).
		Padding(0, 1)

	t.ChipBtn = lipgloss.NewStyle().
		Background(lipgloss.Color("#0A121E")).
		Foreground(t.TextMuted).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#1B2E42")).
		Padding(0, 1).
		MarginRight(1)

	t.ChipBtnActive = lipgloss.NewStyle().
		Background(t.Primary).
		Foreground(lipgloss.Color("#000000")).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	t.SectionTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	t.SectionLine = lipgloss.NewStyle().
		Foreground(t.Secondary)

	t.PromptStyle = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)

	t.CommandEcho = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	t.ErrorOutput = lipgloss.NewStyle().
		Foreground(t.Danger).
		Bold(true)

	t.ModalBox = lipgloss.NewStyle().
		Background(t.BgPanel).
		Foreground(t.TextMain).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(t.Primary)
}
