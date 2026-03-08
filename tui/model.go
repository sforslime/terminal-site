package tui

import (
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

type Model struct {
	screen   screen
	navIndex int
	width    int
	height   int
	list     list.Model
}

func NewModel(width, height int) Model {
	return Model{
		screen: screenHome,
		width:  width,
		height: height,
		list:   newProjectList(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
