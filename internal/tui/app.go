package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smford/cyber-hud/internal/config"
	"github.com/smford/cyber-hud/internal/theme"
)

type AppState int

const (
	StateBoot AppState = iota
	StateMain
	StateMatrix
	StateModal
)

type AppModel struct {
	config       *config.Config
	theme        *theme.Theme
	state        AppState
	width        int
	height       int
	mouseEnabled bool

	// Submodels
	bootModel    BootModel
	tabsModel    TabsModel
	overviewTab  OverviewTabModel
	logsTab      LogsTabModel
	firmwareTab  FirmwareTabModel
	operatorTab  OperatorTabModel
	matrixModel  MatrixModel
	dossierModel DossierModel
}

func NewAppModel(cfg *config.Config) AppModel {
	th := theme.Cyberpunk2077()
	if cfg.Settings.DefaultTheme == "altered-carbon" {
		th = theme.AlteredCarbon()
	}

	state := StateBoot
	if cfg.Settings.SkipBoot {
		state = StateMain
	}

	boot := NewBootModel(th)
	tabs := NewTabsModel(th, len(cfg.BuildLogs), len(cfg.Repositories))
	overview := NewOverviewTabModel(cfg, th, 100, 35)
	logs := NewLogsTabModel(cfg, th, 100, 35)
	firmware := NewFirmwareTabModel(cfg, th, 100, 35)
	operator := NewOperatorTabModel(cfg, th, 100, 35)
	matrix := NewMatrixModel(th, 100, 35)
	dossier := NewDossierModel(th, 100, 35)

	return AppModel{
		config:       cfg,
		theme:        th,
		state:        state,
		mouseEnabled: cfg.Settings.MouseEnabled,
		bootModel:    boot,
		tabsModel:    tabs,
		overviewTab:  overview,
		logsTab:      logs,
		firmwareTab:  firmware,
		operatorTab:  operator,
		matrixModel:  matrix,
		dossierModel: dossier,
	}
}

func (m AppModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	if m.state == StateBoot {
		cmds = append(cmds, m.bootModel.Init())
	}
	if m.mouseEnabled {
		cmds = append(cmds, tea.EnableMouseCellMotion)
	}
	return tea.Batch(cmds...)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tabsModel.SetWidth(msg.Width)

		contentH := msg.Height
		if contentH < 10 {
			contentH = 10
		}
		m.overviewTab.SetSize(msg.Width, contentH)
		m.logsTab.SetSize(msg.Width, contentH)
		m.firmwareTab.SetSize(msg.Width, contentH)
		m.operatorTab.SetSize(msg.Width, contentH)
		m.matrixModel = NewMatrixModel(m.theme, msg.Width, msg.Height)
		m.dossierModel.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "f2", "ctrl+m":
			m.mouseEnabled = !m.mouseEnabled
			if m.mouseEnabled {
				cmds = append(cmds, tea.EnableMouseCellMotion)
			} else {
				cmds = append(cmds, tea.DisableMouse)
			}
			return m, tea.Batch(cmds...)

		case "f1":
			m.toggleTheme()
			return m, nil
		}

	case tea.MouseMsg:
		if !m.mouseEnabled {
			return m, nil
		}
		if msg.Action == tea.MouseActionPress {
			// Check theme toggle button on tab bar
			if m.tabsModel.CheckThemeClick(msg.X, msg.Y) {
				m.toggleTheme()
				return m, nil
			}
			// Check tabs click
			if clickedTab, ok := m.tabsModel.CheckMouseClick(msg.X, msg.Y); ok {
				m.tabsModel.SetActive(clickedTab)
				return m, nil
			}
		}
	}

	// State-specific routing
	switch m.state {
	case StateBoot:
		var cmd tea.Cmd
		m.bootModel, cmd = m.bootModel.Update(msg)
		cmds = append(cmds, cmd)
		if m.bootModel.completed {
			m.state = StateMain
		}
		return m, tea.Batch(cmds...)

	case StateMatrix:
		var cmd tea.Cmd
		m.matrixModel, cmd = m.matrixModel.Update(msg)
		cmds = append(cmds, cmd)
		if !m.matrixModel.active {
			m.state = StateMain
		}
		return m, tea.Batch(cmds...)

	case StateModal:
		var cmd tea.Cmd
		m.dossierModel, cmd = m.dossierModel.Update(msg)
		cmds = append(cmds, cmd)
		if !m.dossierModel.active {
			m.state = StateMain
		}
		return m, tea.Batch(cmds...)

	case StateMain:
		// Global Navigation Keys
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			// If overview tab terminal is focused and typing, only intercept tab-switching if not editing
			// Check direct number keys if terminal is not taking raw input or if user pressed Alt/Ctrl
			switch keyMsg.String() {
			case "1":
				if m.tabsModel.Active() != TabOverview {
					m.tabsModel.SetActive(TabOverview)
					return m, nil
				}
			case "2":
				if m.tabsModel.Active() != TabOverview {
					m.tabsModel.SetActive(TabBuildLogs)
					return m, nil
				}
			case "3":
				if m.tabsModel.Active() != TabOverview {
					m.tabsModel.SetActive(TabFirmware)
					return m, nil
				}
			case "4":
				if m.tabsModel.Active() != TabOverview {
					m.tabsModel.SetActive(TabOperator)
					return m, nil
				}
			case "tab":
				m.tabsModel.Next()
				return m, nil
			case "shift+tab":
				m.tabsModel.Prev()
				return m, nil
			}
		}

		// Route to active tab
		switch m.tabsModel.Active() {
		case TabOverview:
			var cmd tea.Cmd
			m.overviewTab.terminal, cmd = m.overviewTab.terminal.Update(msg)
			cmds = append(cmds, cmd)

			// Check if terminal triggered rain
			if m.overviewTab.terminal.triggerRain {
				m.overviewTab.terminal.triggerRain = false
				m.state = StateMatrix
				m.matrixModel.active = true
				cmds = append(cmds, m.matrixModel.Init())
				return m, tea.Batch(cmds...)
			}
			// Check if terminal triggered dossier
			if m.overviewTab.terminal.triggerModal != "" {
				id := m.overviewTab.terminal.triggerModal
				m.overviewTab.terminal.triggerModal = ""
				for _, log := range m.config.BuildLogs {
					if strings.EqualFold(log.ID, id) {
						m.dossierModel.Open(&log)
						m.state = StateModal
						break
					}
				}
			}
			// Check reboot
			if m.overviewTab.terminal.triggerBoot {
				m.overviewTab.terminal.triggerBoot = false
				m.bootModel = NewBootModel(m.theme)
				m.state = StateBoot
				cmds = append(cmds, m.bootModel.Init())
				return m, tea.Batch(cmds...)
			}

		case TabBuildLogs:
			var cmd tea.Cmd
			m.logsTab, cmd = m.logsTab.Update(msg)
			cmds = append(cmds, cmd)
			if m.logsTab.triggerModal != "" {
				id := m.logsTab.triggerModal
				m.logsTab.triggerModal = ""
				for _, log := range m.config.BuildLogs {
					if strings.EqualFold(log.ID, id) {
						m.dossierModel.Open(&log)
						m.state = StateModal
						break
					}
				}
			}

		case TabFirmware:
			var cmd tea.Cmd
			m.firmwareTab, cmd = m.firmwareTab.Update(msg)
			cmds = append(cmds, cmd)

		case TabOperator:
			var cmd tea.Cmd
			m.operatorTab, cmd = m.operatorTab.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *AppModel) toggleTheme() {
	if m.theme.ID == "cyberpunk" {
		m.theme = theme.AlteredCarbon()
	} else {
		m.theme = theme.Cyberpunk2077()
	}
	m.tabsModel.SetTheme(m.theme)
	m.overviewTab.SetTheme(m.theme)
	m.logsTab.SetTheme(m.theme)
	m.firmwareTab.SetTheme(m.theme)
	m.operatorTab.SetTheme(m.theme)
	m.dossierModel.SetTheme(m.theme)
	m.matrixModel.theme = m.theme
}

func (m AppModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing Cyberdeck..."
	}

	switch m.state {
	case StateBoot:
		return m.bootModel.View()

	case StateMatrix:
		return m.matrixModel.View()

	case StateModal:
		modal := m.dossierModel.View()
		// Overlay modal in center
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(m.theme.BgBase),
		)

	case StateMain:
		return m.renderMainHUD()
	}

	return ""
}

func (m AppModel) renderMainHUD() string {
	var content string

	// Active Tab Content
	switch m.tabsModel.Active() {
	case TabOverview:
		content = m.overviewTab.View()
	case TabBuildLogs:
		content = m.logsTab.View()
	case TabFirmware:
		content = m.firmwareTab.View()
	case TabOperator:
		content = m.operatorTab.View()
	}

	// SRE Invariant: Ensure lines never exceed terminal height to prevent vertical scroll displacement
	lines := strings.Split(content, "\n")
	if m.height > 0 && len(lines) > m.height {
		lines = lines[:m.height]
		content = strings.Join(lines, "\n")
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(m.theme.BgBase).
		Render(content)
}
