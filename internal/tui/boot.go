package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smford/cyber-hud/internal/theme"
)

type BootStep struct {
	Text  string
	Kind  string // "stage", "ok", "info", "warn"
	Delay time.Duration
}

var defaultBootSteps = []BootStep{
	{Text: "ROM CHECK: 0x00000000 -> 0x7FFFFFFF OK [32768MB ECC RAM VERIFIED]", Kind: "stage", Delay: 40 * time.Millisecond},
	{Text: "[  OK  ] Initialized Xtensa LX6 Dual-Core SMP processor @ 240MHz", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[ INFO ] Loading Linux 6.1.0-debian-smford (x86_64/arm64 unified ABI)", Kind: "info", Delay: 40 * time.Millisecond},
	{Text: "[  OK  ] Mounted Virtual File System (/dev/vfs) with LittleFS & SPIFFS hooks", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Loaded smford/esp32-asyncwebserver-fileupload-example (104★ / 30 Forks)", Kind: "ok", Delay: 40 * time.Millisecond},
	{Text: "[  OK  ] Initialized smford/raspberry-pi-streaming-camera RTSP pipeline (83★ / 5 Forks)", Kind: "ok", Delay: 40 * time.Millisecond},
	{Text: "[ INFO ] Calibrating 13.56MHz RC522 SPI RFID reader bus on GPIO_05/18/23", Kind: "info", Delay: 45 * time.Millisecond},
	{Text: "[  OK  ] Safety interlock relay isolation armed: smford/eeh-esp32-rfid", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[  OK  ] Mounted SNMP daemon configuration: smford/rpi-snmpd-configuration", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[ INFO ] Polling MikroTik RouterOS v7 DHCP leases via smford/dhcp-routeros-tool", Kind: "info", Delay: 45 * time.Millisecond},
	{Text: "[  OK  ] Mounted Sunton ESP32 touchscreen display CAD: smford/Sunton-ESP32-2432S028R-Case", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[  OK  ] Mounted workshop safety CAD enclosures: smford/emergency-stop-holder", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Initialized temperature telemetry regulator: smford/esp32-temp-monitor", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[  OK  ] Starting event-driven surveillance API daemon: smford/motioneye-snapshotter", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[  OK  ] Initialized Schwarzschild geodesic raymarcher: smford/blackhole", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[  OK  ] Loaded Selenographic LOD reconnaissance engine: smford/fly-around-the-moon", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[  OK  ] Mounted terminal Markdown parser pipeline: smford/mdee", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "[  OK  ] Armed terminal separator utility: smford/delim", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Synced Homebrew package tap formulas: smford/homebrew-tap", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Initialized VT100 Matrix digital rain streaming pipeline: smford/matrix-rain", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Mounted Matrix cascading text viewer: smford/matrix-cat", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Initialized power tool receipt optimization engine: smford/toolsplit", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Loaded in-browser Kubernetes RBAC security analyzer: smford/rbac-editor", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Loaded in-browser YAML AST & anchor validator: smford/yaml-web-validator", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Dialed PCBoard BBS node emulator via CP437 canvas: smford/pcboard-web-emu", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Initialized terminal Mermaid diagram rendering engine: smford/golang-mermaid", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Loaded IPv4/IPv6 CIDR decomposition & route summarizer: smford/cidr-calculator", Kind: "ok", Delay: 30 * time.Millisecond},
	{Text: "[ WARN ] Unrestricted ICE bypass active for local operator 'smford'", Kind: "warn", Delay: 30 * time.Millisecond},
	{Text: "[  OK  ] Started Bash Shell Netrunner Terminal daemon (TTY1)", Kind: "ok", Delay: 35 * time.Millisecond},
	{Text: "SYSTEM: Neural link handshaking complete. Launching Cyberdeck HUD...", Kind: "stage", Delay: 60 * time.Millisecond},
}

type BootTickMsg struct {
	StepIndex int
}

type BootModel struct {
	steps     []BootStep
	current   int
	completed bool
	theme     *theme.Theme
	width     int
	height    int
}

func NewBootModel(th *theme.Theme) BootModel {
	return BootModel{
		steps:     defaultBootSteps,
		current:   0,
		completed: false,
		theme:     th,
	}
}

func (m BootModel) Init() tea.Cmd {
	return m.nextStepCmd(0)
}

func (m BootModel) nextStepCmd(index int) tea.Cmd {
	if index >= len(m.steps) {
		return func() tea.Msg {
			return BootTickMsg{StepIndex: len(m.steps)}
		}
	}
	delay := m.steps[index].Delay
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return BootTickMsg{StepIndex: index + 1}
	})
}

func (m BootModel) Update(msg tea.Msg) (BootModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "enter", "space", "q":
			m.completed = true
			return m, nil
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			m.completed = true
			return m, nil
		}
	case BootTickMsg:
		m.current = msg.StepIndex
		if m.current >= len(m.steps) {
			m.completed = true
			return m, nil
		}
		return m, m.nextStepCmd(m.current)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m BootModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Primary)

	subStyle := lipgloss.NewStyle().
		Foreground(m.theme.TextMuted)

	skipBtnStyle := lipgloss.NewStyle().
		Foreground(m.theme.Secondary).
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.theme.Secondary).
		Padding(0, 1)

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			headerStyle.Render("NEURAL_BIOS v9.2 // SMFORD MAINFRAME"),
			subStyle.Render("Xtensa Dual-Core & Debian Kernel 6.1.0-SMP / 240MHz Low-Level Bus"),
		),
		"   ",
		skipBtnStyle.Render("[ESC] BYPASS_BOOT"),
	)
	b.WriteString(header + "\n\n")

	// Calculate visible log lines based on height
	maxLines := m.height - 10
	if maxLines < 5 {
		maxLines = 15
	}

	start := 0
	if m.current > maxLines {
		start = m.current - maxLines
	}

	for i := start; i < m.current && i < len(m.steps); i++ {
		step := m.steps[i]
		line := step.Text
		switch step.Kind {
		case "ok":
			line = strings.Replace(line, "[  OK  ]", lipgloss.NewStyle().Foreground(m.theme.Success).Bold(true).Render("[  OK  ]"), 1)
		case "info":
			line = strings.Replace(line, "[ INFO ]", lipgloss.NewStyle().Foreground(m.theme.Secondary).Bold(true).Render("[ INFO ]"), 1)
		case "warn":
			line = strings.Replace(line, "[ WARN ]", lipgloss.NewStyle().Foreground(m.theme.Danger).Bold(true).Render("[ WARN ]"), 1)
		case "stage":
			line = lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true).Render(line)
		}
		b.WriteString(line + "\n")
	}

	// Progress bar
	percent := 0
	if len(m.steps) > 0 {
		percent = (m.current * 100) / len(m.steps)
	}
	if percent > 100 {
		percent = 100
	}

	barWidth := m.width - 20
	if barWidth < 20 {
		barWidth = 40
	}
	filledWidth := (percent * barWidth) / 100
	emptyWidth := barWidth - filledWidth

	bar := lipgloss.NewStyle().Foreground(m.theme.Primary).Render(strings.Repeat("█", filledWidth)) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#1a2332")).Render(strings.Repeat("░", emptyWidth))

	b.WriteString("\n" + bar + "\n")
	statusLine := fmt.Sprintf("STAGE: INITIALIZING_POST_ROUTINES...  [%3d%%]", percent)
	b.WriteString(lipgloss.NewStyle().Foreground(m.theme.Secondary).Render(statusLine))

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(1, 2).
		Render(b.String())
}
