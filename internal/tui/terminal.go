package tui

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smford/cyber-hud/internal/config"
	"github.com/smford/cyber-hud/internal/theme"
)

type TerminalModel struct {
	config       *config.Config
	theme        *theme.Theme
	input        textinput.Model
	viewport     viewport.Model
	history      []string
	historyIdx   int
	outputLines  []string
	width        int
	height       int
	chips        []config.QuickChip
	chipFocus    int
	focused      bool
	triggerRain  bool
	triggerModal string // ID of log to open
	triggerBoot  bool
}

func NewTerminalModel(cfg *config.Config, th *theme.Theme, w, h int) TerminalModel {
	ti := textinput.New()
	ti.Placeholder = "enter command ('help', 'matrix-rain', 'cidr', 'projects', 'cat notes.md')..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = w - 24

	vp := viewport.New(w, h-10)
	initialOutput := []string{
		"SYSTEM: Neural Debian Kernel 6.1.0-debian-smford loaded.",
		"SYSTEM: Cyberdeck CLI active across all 27 projects (Type 'help' or click quick chips).",
		"SYSTEM: Views: [1] OVERVIEW  [2] BUILD_LOGS  [3] FIRMWARE_DECK  [4] OPERATOR_STACK (Press 1-4)",
		"SYSTEM: Try 'matrix-rain', 'cidr', 'projects', 'blackhole', 'moon', 'cat notes.md'.",
	}

	m := TerminalModel{
		config:      cfg,
		theme:       th,
		input:       ti,
		viewport:    vp,
		history:     []string{},
		historyIdx:  -1,
		outputLines: initialOutput,
		width:       w,
		height:      h,
		chips:       cfg.QuickChips,
		chipFocus:   -1,
		focused:     true,
	}
	m.refreshViewport()
	return m
}

func (m *TerminalModel) SetTheme(th *theme.Theme) {
	m.theme = th
}

func (m *TerminalModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.input.Width = w - 24
	if m.input.Width < 20 {
		m.input.Width = 20
	}
	vpH := h - 9
	if vpH < 5 {
		vpH = 5
	}
	m.viewport.Width = w - 4
	m.viewport.Height = vpH
	m.refreshViewport()
}

func (m *TerminalModel) Focus() {
	m.focused = true
	m.input.Focus()
}

func (m *TerminalModel) Blur() {
	m.focused = false
	m.input.Blur()
}

func (m *TerminalModel) refreshViewport() {
	content := strings.Join(m.outputLines, "\n")
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

func (m *TerminalModel) AppendOutput(line string) {
	m.outputLines = append(m.outputLines, line)
	m.refreshViewport()
}

func (m TerminalModel) Update(msg tea.Msg) (TerminalModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			rawCmd := strings.TrimSpace(m.input.Value())
			if rawCmd != "" {
				m.history = append(m.history, rawCmd)
				m.historyIdx = len(m.history)
				m.input.SetValue("")

				promptEcho := fmt.Sprintf("%s@%s:~$ %s", m.config.Operator.Handle, m.config.Operator.Host, rawCmd)
				m.outputLines = append(m.outputLines, promptEcho)
				m.execute(rawCmd)
				m.refreshViewport()
			}
			return m, nil

		case "up":
			if len(m.history) > 0 {
				if m.historyIdx > 0 {
					m.historyIdx--
					m.input.SetValue(m.history[m.historyIdx])
					m.input.CursorEnd()
				}
			}
			return m, nil

		case "down":
			if len(m.history) > 0 {
				if m.historyIdx < len(m.history)-1 {
					m.historyIdx++
					m.input.SetValue(m.history[m.historyIdx])
					m.input.CursorEnd()
				} else {
					m.historyIdx = len(m.history)
					m.input.SetValue("")
				}
			}
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

	if m.focused {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m *TerminalModel) execute(cmdStr string) {
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return
	}
	parts := strings.Fields(cmdStr)
	mainCmd := strings.ToLower(parts[0])

	switch mainCmd {
	case "clear":
		m.outputLines = []string{}
		return

	case "help":
		m.AppendOutput(m.helpText())
		return

	case "projects":
		m.AppendOutput(m.projectsList())
		return

	case "matrix-rain":
		m.triggerRain = true
		return

	case "reboot", "boot":
		m.triggerBoot = true
		return

	case "ls":
		m.AppendOutput(m.listDirectory())
		return

	case "cat":
		if len(parts) < 2 {
			m.AppendOutput(m.theme.ErrorOutput.Render("cat: missing filename. Usage: cat <file> (Try 'ls' for files)"))
			return
		}
		filename := parts[1]
		if content, ok := m.config.VFS[filename]; ok {
			m.AppendOutput(content)
		} else {
			m.AppendOutput(m.theme.ErrorOutput.Render(fmt.Sprintf("cat: %s: No such file or directory. Available: bio.txt, hardware.ino, temp_monitor.ino, snmpd.conf, notes.md", filename)))
		}
		return

	case "which":
		if len(parts) < 2 {
			m.AppendOutput(m.theme.ErrorOutput.Render("which: missing operand"))
			return
		}
		target := parts[1]
		bins := []string{
			"matrix-rain", "cidr", "mcat", "mdee", "delim", "projects", "toolsplit",
			"blackhole", "moon", "bbs", "burger", "rbac", "yaml", "snmp", "rfid", "cam",
		}
		found := false
		for _, b := range bins {
			if b == target {
				found = true
				break
			}
		}
		if found {
			m.AppendOutput(fmt.Sprintf("/usr/local/bin/%s", target))
		} else {
			m.AppendOutput(m.theme.ErrorOutput.Render(fmt.Sprintf("which: no %s in (/usr/local/bin:/usr/bin:/bin)", target)))
		}
		return

	case "cidr":
		if len(parts) < 2 {
			m.AppendOutput("Usage: cidr <IP/CIDR> or cidr <IP>+<count> (e.g. 'cidr 192.168.1.0/24' or 'cidr 192.168.1.10+5')")
			return
		}
		m.AppendOutput(m.executeCidr(parts[1]))
		return

	case "delim":
		banner := strings.Repeat("━", 48)
		m.AppendOutput(m.theme.SectionLine.Render(banner))
		m.AppendOutput("[DELIM] ZERO-ALLOCATION TERMINAL STREAM BARRIER ACTIVE")
		m.AppendOutput(m.theme.SectionLine.Render(banner))
		return

	case "mcat":
		if len(parts) < 2 {
			m.AppendOutput("mcat: Cascading matrix digital rain source viewer. Usage: mcat <file> (e.g. 'mcat notes.md')")
			return
		}
		filename := parts[1]
		if content, ok := m.config.VFS[filename]; ok {
			m.AppendOutput(fmt.Sprintf("Cascading source stream for [%s]:\n%s", filename, content))
		} else {
			m.AppendOutput(m.theme.ErrorOutput.Render(fmt.Sprintf("mcat: %s: file not found in VFS", filename)))
		}
		return

	case "open":
		if len(parts) < 2 {
			m.AppendOutput("open: missing target. e.g. 'open blackhole', 'open moon', 'open cidr'")
			return
		}
		target := parts[1]
		// Search repositories
		for _, repo := range m.config.Repositories {
			if strings.EqualFold(repo.Name, target) {
				url := repo.URL
				if repo.Demo != "" {
					url = repo.Demo
				}
				openBrowser(url)
				m.AppendOutput(fmt.Sprintf("Launching [%s] in browser: %s", repo.Name, url))
				return
			}
		}
		// Search custom commands
		if cmd, ok := m.config.CustomCommands[target]; ok && cmd.URL != "" {
			openBrowser(cmd.URL)
			m.AppendOutput(fmt.Sprintf("Launching [%s] in browser: %s", target, cmd.URL))
			return
		}
		m.AppendOutput(m.theme.ErrorOutput.Render(fmt.Sprintf("open: target '%s' not found", target)))
		return

	case "view":
		if len(parts) < 2 {
			m.AppendOutput("view: missing log ID. e.g. 'view log-042'")
			return
		}
		id := strings.ToLower(parts[1])
		for _, log := range m.config.BuildLogs {
			if strings.EqualFold(log.ID, id) {
				m.triggerModal = log.ID
				m.AppendOutput(fmt.Sprintf("Accessing schematic dossier for [%s]: %s", log.ID, log.Title))
				return
			}
		}
		m.AppendOutput(m.theme.ErrorOutput.Render(fmt.Sprintf("view: build log '%s' not found", id)))
		return

	default:
		// Check custom commands
		if c, ok := m.config.CustomCommands[mainCmd]; ok {
			if len(parts) > 1 && (parts[1] == "web" || parts[1] == "open") && c.URL != "" {
				openBrowser(c.URL)
				m.AppendOutput(fmt.Sprintf("Launching [%s] in browser: %s", mainCmd, c.URL))
			} else {
				m.AppendOutput(c.Output)
			}
			return
		}

		// Check repositories match
		for _, repo := range m.config.Repositories {
			if strings.EqualFold(repo.Name, mainCmd) {
				m.AppendOutput(fmt.Sprintf("[%s] ★ %d | Forks %d | %s\nTech: %s\nURL: %s",
					repo.Name, repo.Stars, repo.Forks, repo.Category, strings.Join(repo.Tech, ", "), repo.URL))
				return
			}
		}

		m.AppendOutput(m.theme.ErrorOutput.Render(fmt.Sprintf("bash: %s: command not found. Type 'help' for directory.", mainCmd)))
	}
}

func (m TerminalModel) helpText() string {
	return `MAINFRAME COMMAND PROTOCOL DIRECTORY:

SYSTEM & HUD:
  • projects       - Dump all 27 GitHub repositories with telemetry & stars
  • cd <1-4>       - Switch HUD view [1:overview | 2:logs | 3:firmware | 4:operator]
  • open <project> - Launch any web app or GitHub repo (e.g. 'open blackhole')
  • matrix-rain    - Launch authentic VT100 Matrix digital rain screensaver
  • mcat <file>    - Cascading digital rain text & code viewer (e.g. 'mcat notes.md')
  • cidr <range>   - High-perf IPv4/IPv6 CIDR decomposition & collision calculator
  • delim          - Zero-alloc terminal visual barrier pipeline
  • view <log-id>  - Open schematic dossier modal (e.g. 'view log-042')
  • ls / which     - List executable binaries and Virtual File System (VFS)
  • cat <file>     - Display raw file contents (notes.md, bio.txt, hardware.ino)
  • clear / reboot - Clear terminal buffer / replay hardware POST boot
  • theme          - Toggle Night City / Altered Carbon cyberware styling

INTERACTIVE TELEMETRY & DISPATCHES:
  • blackhole, moon, bbs, burger, toolsplit, rbac, yaml, esp32, cam, rfid, temp, cad, estop`
}

func (m TerminalModel) projectsList() string {
	var sb strings.Builder
	sb.WriteString("GITHUB REPOSITORY INDEX (27 NODES):\n")
	for i, r := range m.config.Repositories {
		sb.WriteString(fmt.Sprintf("  [%02d] %-42s (★ %3d) [%s]\n", i+1, r.Name, r.Stars, strings.Join(r.Tech, "/")))
	}
	return sb.String()
}

func (m TerminalModel) listDirectory() string {
	bins := []string{
		"-rwxr-xr-x 1 smford staff  4.5M  matrix-rain*",
		"-rwxr-xr-x 1 smford staff  4.1M  cidr*",
		"-rwxr-xr-x 1 smford staff  3.8M  mdee*",
		"-rwxr-xr-x 1 smford staff  3.4M  mcat*",
		"-rwxr-xr-x 1 smford staff  3.2M  delim*",
		"-rwxr-xr-x 1 smford staff  2.9M  toolsplit*",
		"-rwxr-xr-x 1 smford staff  2.7M  blackhole*",
		"-rwxr-xr-x 1 smford staff  2.5M  moon*",
	}
	var vfsFiles []string
	for name, content := range m.config.VFS {
		vfsFiles = append(vfsFiles, fmt.Sprintf("-rw-r--r-- 1 smford staff  %5dB  %s", len(content), name))
	}
	total := len(bins) + len(vfsFiles)
	return fmt.Sprintf("total %d\n%s\n%s", total, strings.Join(bins, "\n"), strings.Join(vfsFiles, "\n"))
}

func (m TerminalModel) executeCidr(arg string) string {
	if strings.Contains(arg, "+") {
		parts := strings.Split(arg, "+")
		baseIP := net.ParseIP(parts[0])
		count, _ := strconv.Atoi(parts[1])
		if baseIP == nil || count <= 0 {
			return "Invalid CIDR addition format. Example: 192.168.1.10+5"
		}
		return fmt.Sprintf("[CIDR DECOMPOSITION]\nBase IP: %s | Host Count: %d\nRange: %s - %s",
			baseIP.String(), count, baseIP.String(), net.IPv4(baseIP[12], baseIP[13], baseIP[14], baseIP[15]+byte(count)).String())
	}

	ip, ipnet, err := net.ParseCIDR(arg)
	if err != nil {
		return fmt.Sprintf("Error parsing CIDR: %v", err)
	}
	mask := ipnet.Mask
	ones, bits := mask.Size()
	return fmt.Sprintf("[CIDR REPORT: %s]\nIP: %s | Mask: /%d (%s) | Total Addresses: %d",
		arg, ip.String(), ones, net.IP(mask).String(), 1<<(bits-ones))
}

func (m TerminalModel) View() string {
	var sb strings.Builder

	// Terminal Header
	dots := lipgloss.NewStyle().Foreground(m.theme.Danger).Render("●") + " " +
		lipgloss.NewStyle().Foreground(m.theme.Primary).Render("●") + " " +
		lipgloss.NewStyle().Foreground(m.theme.Secondary).Render("●")

	leftHeader := lipgloss.JoinHorizontal(
		lipgloss.Center,
		dots,
		"  ",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render("BASH NETRUNNER CYBERDECK v4.8"),
		"  ",
		lipgloss.NewStyle().Foreground(m.theme.Primary).Render("[TTY1_ONLINE]"),
	)

	viewsBadge := lipgloss.NewStyle().Foreground(m.theme.TextMuted).Render("[1:HUD 2:LOGS 3:DECK 4:STACK]")
	styleBadge := lipgloss.NewStyle().Foreground(m.theme.Secondary).Render(fmt.Sprintf("[%s (F1)]", m.theme.Name))
	rightHeader := lipgloss.JoinHorizontal(lipgloss.Center, viewsBadge, "  ", styleBadge)

	gap := (m.width - 6) - lipgloss.Width(leftHeader) - lipgloss.Width(rightHeader)
	var termHeader string
	if gap > 2 {
		termHeader = leftHeader + strings.Repeat(" ", gap) + rightHeader
	} else {
		termHeader = leftHeader
	}
	sb.WriteString(termHeader + "\n")

	// Viewport (Output)
	sb.WriteString(m.viewport.View() + "\n")

	// Quick Exec Chips
	var renderedChips []string
	renderedChips = append(renderedChips, lipgloss.NewStyle().Foreground(m.theme.TextMuted).Render("Quick:"))
	for _, chip := range m.chips[:10] { // Top 10 chips
		renderedChips = append(renderedChips, m.theme.ChipBtn.Render(chip.Label))
	}
	chipsLine := lipgloss.JoinHorizontal(lipgloss.Center, renderedChips...)
	sb.WriteString(chipsLine + "\n")

	// Input Line
	prompt := m.theme.PromptStyle.Render(fmt.Sprintf("%s@%s:~$ ", m.config.Operator.Handle, m.config.Operator.Host))
	inputLine := lipgloss.JoinHorizontal(lipgloss.Center, prompt, m.input.View())
	sb.WriteString(inputLine)

	return lipgloss.NewStyle().
		Width(m.width - 2).
		Background(m.theme.BgPanel).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Secondary).
		Padding(0, 1).
		Render(sb.String())
}

// CheckChipClick tests if mouse clicked on one of the quick chips.
func (m *TerminalModel) CheckChipClick(x, y, termStartY int) (string, bool) {
	// Chips line is located near bottom of terminal window
	chipsY := termStartY + m.height - 4
	if y != chipsY {
		return "", false
	}
	currentX := 8 // after "Quick:"
	for _, chip := range m.chips[:10] {
		chipW := len(chip.Label) + 4
		if x >= currentX && x < currentX+chipW {
			return chip.Command, true
		}
		currentX += chipW + 1
	}
	return "", false
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "darwin":
		_ = exec.Command("open", url).Start()
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	}
}
