package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var snowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

func renderLogoWithSnow(snow []snowflake) string {
	logoLines := strings.Split(nameLogo, "\n")

	// build snow lookup: (col, row) -> char
	flakeMap := make(map[[2]int]rune, len(snow))
	for _, f := range snow {
		flakeMap[[2]int{f.x, f.y}] = f.ch
	}

	var sb strings.Builder
	for row, line := range logoLines {
		runes := []rune(line)
		for col := 0; col < snowW; col++ {
			var ch rune = ' '
			if col < len(runes) {
				ch = runes[col]
			}
			if ch != ' ' {
				sb.WriteString(accentStyle.Render(string(ch)))
			} else if flake, ok := flakeMap[[2]int{col, row}]; ok {
				sb.WriteString(snowStyle.Render(string(flake)))
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

var (
	accent    = lipgloss.Color("#4ecdc4")
	subtle    = lipgloss.Color("#444444")
	dim       = lipgloss.Color("#2a2a2a")
	white     = lipgloss.Color("#e8e8e8")
	whiteHigh = lipgloss.Color("#ffffff")

	accentStyle  = lipgloss.NewStyle().Foreground(accent)
	subtleStyle  = lipgloss.NewStyle().Foreground(subtle)
	dimStyle     = lipgloss.NewStyle().Foreground(dim)
	bodyStyle    = lipgloss.NewStyle().Foreground(white)
	boldStyle    = lipgloss.NewStyle().Bold(true).Foreground(whiteHigh)
	footerStyle  = lipgloss.NewStyle().Foreground(subtle)

	navActive   = lipgloss.NewStyle().Foreground(accent).Bold(true)
	navInactive = lipgloss.NewStyle().Foreground(white)
)

const nameLogo = `
  :::.  .-:.     ::-.   ...      .:
  ;;` + "`" + `;;  ';;.   ;;;;'.;;;;;;;.  ;;;
 ,[[ '[[,  '[[,[[[' ,[[     \[[,'[[
c$$$cc$$$c   c$$"   $$$,     $$$ $$
 888   888,,8P"` + "`" + `    "888,_ _,88P ""
 YMM   "` + "`" + `mM"         "YMMMMMP"  MM
`

func viewHome(m Model) string {
	logo := renderLogoWithSnow(m.snow)

	// Bio paragraphs
	bioPrimary := boldStyle.Render(
		"is a creator & storyteller on the internet,\nbuilding cool products,\ndocumenting life & reflecting on how\ntechnology shapes our humanity.",
	)

	bioSecondary := "\n\n" + boldStyle.Render(
		"AYO also works as a builder & maker,\ncreating things people actually use.",
	)

	bioFaded := "\n\n" + subtleStyle.Render(
		"Rooted in curiosity, obsessed with craft.\nWork sits at the intersection of\nhuman nature, the arts, and technology.",
	)

	explore := "\n\n" + dimStyle.Render("Explore the directories below ↓")

	// Horizontal nav
	nav := buildNav(m.navIndex)

	content := logo +
		bioPrimary + bioSecondary + bioFaded + explore +
		"\n\n" + nav

	main := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 1).
		Padding(1, 4).
		Render(content)

	footer := footerStyle.Width(m.width).Render(
		"[← → to select · enter to open · q to quit]",
	)

	return lipgloss.JoinVertical(lipgloss.Left, main, footer)
}

func buildNav(selected int) string {
	var parts []string
	for i, item := range navItems {
		if i == selected {
			parts = append(parts, navActive.Render("+ "+item))
		} else {
			parts = append(parts, navInactive.Render(item))
		}
	}
	return strings.Join(parts, "   ")
}

func viewCreations(m Model) string {
	header := boldStyle.Render("Creations") + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	body := m.list.View()
	footer := footerStyle.Render("\nesc go back • q quit")

	return lipgloss.NewStyle().Padding(1, 2).Render(header + body + footer)
}

func viewReflections(_ Model) string {
	header := boldStyle.Render("Reflections") + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	thoughts := []string{
		"→ building in public is underrated",
		"→ terminals are timeless",
		"→ simplicity beats complexity, always",
	}

	body := bodyStyle.Render(strings.Join(thoughts, "\n"))
	footer := footerStyle.Render("\nesc go back • q quit")

	return lipgloss.NewStyle().Padding(1, 2).Render(header + body + footer)
}

func viewContacts(_ Model) string {
	header := boldStyle.Render("Contacts") + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	links := fmt.Sprintf(
		"%s  https://www.github.com/sforslime\n%s  https://www.x.com/sforslime\n%s  https://www.instagram.com/yourstruly.ayo/",
		accentStyle.Render("gh"),
		accentStyle.Render("tw"),
		accentStyle.Render("@"),
	)

	footer := footerStyle.Render("\nesc go back • q quit")

	return lipgloss.NewStyle().Padding(1, 2).Render(header + links + footer)
}

// project list

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
