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

	accentStyle = lipgloss.NewStyle().Foreground(accent)
	subtleStyle = lipgloss.NewStyle().Foreground(subtle)
	dimStyle    = lipgloss.NewStyle().Foreground(dim)
	bodyStyle   = lipgloss.NewStyle().Foreground(white)
	boldStyle   = lipgloss.NewStyle().Bold(true).Foreground(whiteHigh)
	footerStyle = lipgloss.NewStyle().Foreground(subtle)

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
	const portraitWidth = 62 // 60 chars + 2 padding
	rightWidth := m.width - portraitWidth

	// --- Left: animated portrait (fixed clip at 26 lines) ---
	portraitLines := strings.Split(m.frames[m.frameIndex], "\n")
	if len(portraitLines) > 26 {
		portraitLines = portraitLines[:26]
	}
	left := lipgloss.NewStyle().
		Width(portraitWidth).
		Height(m.height-1).
		Padding(0, 1).
		Render(strings.Join(portraitLines, "\n"))

	// --- Right: logo + bio + nav ---
	logo := renderLogoWithSnow(m.snow)

	bioPrimary := boldStyle.Render(
		"is a computer science major & creator on the internet,\nbuilding cool software, documenting life & reflecting\non how technology shapes our humanity.",
	)

	bioSecondary := "\n\n" + boldStyle.Render(
		"AYO also works as a developer & engineer,\ncreating things people actually use.",
	)

	bioFaded := "\n\n" + subtleStyle.Render(
		"Rooted in curiosity, obsessed with craft.\nWork sits at the intersection of\nhuman nature, the arts, and technology.",
	)

	explore := "\n\n" + dimStyle.Render("Explore the directories below ↓")

	nav := buildNav(m.navIndex)

	rightContent := logo +
		bioPrimary + bioSecondary + bioFaded + explore +
		"\n\n" + nav

	right := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height-1).
		Padding(1, 2).
		Render(rightContent)

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	footer := footerStyle.Width(m.width).Render(
		"[← → to select · enter to open · q to quit]",
	)

	return lipgloss.JoinVertical(lipgloss.Left, body, footer)
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


func viewContacts(m Model) string {
	header := boldStyle.Render("Contacts") + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	links := fmt.Sprintf(
		"%s  github.com/sforslime\n%s  x.com/sforslime\n%s  instagram.com/yourstruly.ayo",
		accentStyle.Render("Github   "),
		accentStyle.Render("X        "),
		accentStyle.Render("Instagram"),
	)

	footer := footerStyle.Render("\nesc go back • q quit")

	starLines := strings.Split(m.star, "\n")
	// show the middle section where the star's arms are visible
	start := 8
	end := 24
	if end > len(starLines) {
		end = len(starLines)
	}
	if start > len(starLines) {
		start = 0
	}
	star := accentStyle.Render(strings.Join(starLines[start:end], "\n"))

	left := lipgloss.NewStyle().Width(36).Render(header + links + footer)
	right := lipgloss.NewStyle().Padding(1, 8).Render(star)

	return lipgloss.NewStyle().Padding(1, 2).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
	)
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
		project{"web scraper", "a scraper - built with Selenium + Python"},
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
