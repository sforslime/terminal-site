package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	accent    = lipgloss.Color("#4ecdc4")
	subtle    = lipgloss.Color("#555555")
	highlight = lipgloss.Color("#ffffff")

	accentStyle = lipgloss.NewStyle().Foreground(accent)
	subtleStyle = lipgloss.NewStyle().Foreground(subtle)
	boldStyle   = lipgloss.NewStyle().Bold(true).Foreground(highlight)

	navNormal   = lipgloss.NewStyle().Foreground(subtle).PaddingLeft(2)
	navSelected = lipgloss.NewStyle().Foreground(accent).Bold(true).PaddingLeft(1).SetString("> ")

	footerStyle = lipgloss.NewStyle().Foreground(subtle).MarginTop(1)
	panelStyle  = lipgloss.NewStyle().Padding(1, 2)
)

const asciiArt = `
   ___  __ ____  __
  / _ |/ // __ \/ /
 / __ |/ // /_/ /_/
/_/ |_/_/ \____(_)
`

func viewHome(m Model) string {
	// Left panel: ASCII art
	left := panelStyle.Copy().Width(m.width / 2).Render(
		accentStyle.Render(asciiArt) + "\n" +
			subtleStyle.Render("  ssh portfolio"),
	)

	// Right panel: bio + nav
	bio := boldStyle.Render("AYO!") + "\n" +
		subtleStyle.Render("developer. builder. maker of things.") + "\n\n"

	nav := ""
	for i, item := range navItems {
		if i == m.navIndex {
			nav += navSelected.Render(item) + "\n"
		} else {
			nav += navNormal.Render(item) + "\n"
		}
	}

	footer := footerStyle.Render("↑/↓ navigate • enter select • q quit")

	right := panelStyle.Copy().Width(m.width / 2).Render(
		bio + nav + "\n" + footer,
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func viewCreations(m Model) string {
	header := boldStyle.Render("Creations") + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	body := m.list.View()
	footer := footerStyle.Render("\nesc go back • q quit")

	return panelStyle.Render(header + body + footer)
}

func viewReflections(m Model) string {
	header := boldStyle.Render("Reflections") + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	thoughts := []string{
		"→ building in public is underrated",
		"→ terminals are timeless",
		"→ simplicity beats complexity, always",
	}

	body := subtleStyle.Render(strings.Join(thoughts, "\n"))
	footer := footerStyle.Render("\nesc go back • q quit")

	return panelStyle.Render(header + body + footer)
}

func viewContacts(m Model) string {
	header := boldStyle.Render("Contacts") + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	links := fmt.Sprintf(
		"%s  github.com/sforslime\n%s  add your twitter/x\n%s  add your email",
		accentStyle.Render("gh"),
		accentStyle.Render("tw"),
		accentStyle.Render("@"),
	)

	footer := footerStyle.Render("\nesc go back • q quit")

	return panelStyle.Render(header + links + footer)
}

// project list setup

type project struct {
	title, desc string
}

func (p project) Title() string       { return p.title }
func (p project) Description() string { return p.desc }
func (p project) FilterValue() string { return p.title }

func newProjectList() list.Model {
	projects := []list.Item{
		project{"ssh-portfolio", "this terminal — built with Bubble Tea + Wish"},
		project{"project two", "add your projects here"},
		project{"project three", "keep building"},
	}

	l := list.New(projects, list.NewDefaultDelegate(), 60, 10)
	l.Title = ""
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle()
	return l
}
