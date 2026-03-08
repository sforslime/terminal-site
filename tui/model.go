package tui

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenHome screen = iota
	screenCreations
	screenReflections
	screenContacts
)

var navItems = []string{"Creations", "Reflections", "Contacts"}

const snowW = 48
const snowH = 10

type snowflake struct {
	x, y int
	ch   rune
}

var snowChars = []rune{'*', '.', '\'', ',', '+', '·'}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initSnow() []snowflake {
	flakes := make([]snowflake, 12)
	for i := range flakes {
		flakes[i] = snowflake{
			x:  rand.Intn(snowW),
			y:  rand.Intn(snowH),
			ch: snowChars[rand.Intn(len(snowChars))],
		}
	}
	return flakes
}

func tickSnow(snow []snowflake) []snowflake {
	next := snow[:0]
	for _, f := range snow {
		f.y++
		if f.y < snowH {
			next = append(next, f)
		}
	}
	count := rand.Intn(3) + 1
	for i := 0; i < count; i++ {
		next = append(next, snowflake{
			x:  rand.Intn(snowW),
			y:  0,
			ch: snowChars[rand.Intn(len(snowChars))],
		})
	}
	return next
}

type Model struct {
	screen   screen
	navIndex int
	width    int
	height   int
	list     list.Model
	snow     []snowflake
}

func NewModel(width, height int) Model {
	return Model{
		screen: screenHome,
		width:  width,
		height: height,
		list:   newProjectList(),
		snow:   initSnow(),
	}
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.snow = tickSnow(m.snow)
		return m, doTick()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "left", "h":
			if m.screen == screenHome {
				if m.navIndex > 0 {
					m.navIndex--
				}
			}

		case "right", "l":
			if m.screen == screenHome {
				if m.navIndex < len(navItems)-1 {
					m.navIndex++
				}
			}

		case "up", "k":
			if m.screen == screenCreations {
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}

		case "down", "j":
			if m.screen == screenCreations {
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}

		case "enter":
			if m.screen == screenHome {
				m.screen = screen(m.navIndex + 1)
			}

		case "esc", "backspace":
			m.screen = screenHome
		}
	}

	if m.screen == screenCreations {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case screenCreations:
		return viewCreations(m)
	case screenReflections:
		return viewReflections(m)
	case screenContacts:
		return viewContacts(m)
	default:
		return viewHome(m)
	}
}
