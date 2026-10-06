package tui

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smford/cyber-hud/internal/theme"
)

type MatrixTickMsg time.Time

type MatrixColumn struct {
	HeadY  int
	Length int
	Speed  int
	Chars  []rune
}

type MatrixModel struct {
	theme     *theme.Theme
	width     int
	height    int
	columns   []MatrixColumn
	active    bool
	katakana  []rune
	frameRate time.Duration
}

func NewMatrixModel(th *theme.Theme, w, h int) MatrixModel {
	chars := []rune("ｦｱｳｴｵｶｷｹｺｻｼｽｾｿﾀﾂﾃﾅﾆﾇﾈﾊﾋﾎﾏﾐﾑﾒﾓﾔﾕﾗﾘﾜ0123456789ABCDEF$#@*+=-")
	m := MatrixModel{
		theme:     th,
		width:     w,
		height:    h,
		active:    false,
		katakana:  chars,
		frameRate: 50 * time.Millisecond,
	}
	m.initColumns()
	return m
}

func (m *MatrixModel) initColumns() {
	colCount := m.width
	if colCount <= 0 {
		colCount = 80
	}
	m.columns = make([]MatrixColumn, colCount)
	for i := range m.columns {
		m.columns[i] = MatrixColumn{
			HeadY:  rand.Intn(m.height + 20) - 20,
			Length: rand.Intn(15) + 8,
			Speed:  rand.Intn(2) + 1,
			Chars:  m.randomChars(30),
		}
	}
}

func (m *MatrixModel) randomChars(count int) []rune {
	r := make([]rune, count)
	for i := range r {
		r[i] = m.katakana[rand.Intn(len(m.katakana))]
	}
	return r
}

func (m MatrixModel) Init() tea.Cmd {
	return tea.Tick(m.frameRate, func(t time.Time) tea.Msg {
		return MatrixTickMsg(t)
	})
}

func (m MatrixModel) Update(msg tea.Msg) (MatrixModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "enter", "ctrl+c":
			m.active = false
			return m, nil
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			m.active = false
			return m, nil
		}
	case MatrixTickMsg:
		if !m.active {
			return m, nil
		}
		// Advance columns
		for i := range m.columns {
			col := &m.columns[i]
			col.HeadY += col.Speed
			if col.HeadY-col.Length > m.height {
				col.HeadY = -rand.Intn(10)
				col.Length = rand.Intn(15) + 8
				col.Speed = rand.Intn(2) + 1
				col.Chars = m.randomChars(30)
			} else if rand.Float32() < 0.2 {
				// Random character mutation
				idx := rand.Intn(len(col.Chars))
				col.Chars[idx] = m.katakana[rand.Intn(len(m.katakana))]
			}
		}
		return m, tea.Tick(m.frameRate, func(t time.Time) tea.Msg {
			return MatrixTickMsg(t)
		})
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initColumns()
	}
	return m, nil
}

func (m MatrixModel) View() string {
	screen := make([][]rune, m.height)
	colorMap := make([][]int, m.height) // 0: none, 1: head (white), 2: bright green, 3: dark green

	for y := 0; y < m.height; y++ {
		screen[y] = make([]rune, m.width)
		colorMap[y] = make([]int, m.width)
		for x := 0; x < m.width; x++ {
			screen[y][x] = ' '
		}
	}

	for x, col := range m.columns {
		if x >= m.width {
			break
		}
		for i := 0; i < col.Length; i++ {
			y := col.HeadY - i
			if y >= 0 && y < m.height {
				charIdx := i % len(col.Chars)
				screen[y][x] = col.Chars[charIdx]
				if i == 0 {
					colorMap[y][x] = 1 // Head
				} else if i < 4 {
					colorMap[y][x] = 2 // Bright
				} else {
					colorMap[y][x] = 3 // Trail
				}
			}
		}
	}

	var sb strings.Builder
	headStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	brightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF66")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#008833"))

	// Render rows except bottom 2 for statusline
	for y := 0; y < m.height-2; y++ {
		for x := 0; x < m.width; x++ {
			r := screen[y][x]
			switch colorMap[y][x] {
			case 1:
				sb.WriteString(headStyle.Render(string(r)))
			case 2:
				sb.WriteString(brightStyle.Render(string(r)))
			case 3:
				sb.WriteString(dimStyle.Render(string(r)))
			default:
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}

	// Statusline
	statusLeft := lipgloss.NewStyle().Foreground(m.theme.Success).Bold(true).Render("MAINFRAME // VT100 MATRIX DIGITAL RAIN STREAM [ACTIVE]")
	statusRight := lipgloss.NewStyle().Foreground(m.theme.Danger).Render("[Q] EXIT STREAM")
	gap := m.width - lipgloss.Width(statusLeft) - lipgloss.Width(statusRight) - 2
	if gap < 0 {
		gap = 2
	}
	statusLine := statusLeft + strings.Repeat(" ", gap) + statusRight

	statusBox := lipgloss.NewStyle().
		Width(m.width).
		Background(lipgloss.Color("#020803")).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(m.theme.Success).
		Padding(0, 1).
		Render(statusLine)

	sb.WriteString(statusBox)
	return sb.String()
}
