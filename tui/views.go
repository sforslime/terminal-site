package tui

import (
	"fmt"
	"strings"

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

	purple      = lipgloss.Color("#a855f7")
	selectStyle = lipgloss.NewStyle().Foreground(purple).Bold(true)

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
		Padding(1, 2).
		Render(rightContent)

	// --- Left: portrait clipped to match right content height ---
	rightLines := strings.Count(right, "\n") + 1
	portraitLines := strings.Split(m.frames[m.frameIndex], "\n")
	if len(portraitLines) > rightLines {
		portraitLines = portraitLines[:rightLines]
	}
	left := lipgloss.NewStyle().
		Width(portraitWidth).
		Padding(0, 1).
		Render(strings.Join(portraitLines, "\n"))

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

	// group by category
	seen := map[string]bool{}
	var body strings.Builder
	for i, c := range allCreations {
		if !seen[c.Category] {
			if len(seen) > 0 {
				body.WriteString("\n")
			}
			body.WriteString(subtleStyle.Render(c.Category) + "\n")
			seen[c.Category] = true
		}
		if i == m.creationIndex {
			body.WriteString(selectStyle.Render("  ✦ "+c.Title) + "\n")
		} else {
			body.WriteString(subtleStyle.Render("    "+c.Title) + "\n")
		}
	}

	footer := footerStyle.Render("\n↑/↓ navigate • enter open • esc back")
	return lipgloss.NewStyle().Padding(1, 2).Render(header + body.String() + footer)
}

func viewCreationDetail(m Model) string {
	c := allCreations[m.creationIndex]

	header := boldStyle.Render(c.Title) + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	category := subtleStyle.Render(c.Category) + "\n\n"
	desc := bodyStyle.Render(c.Desc) + "\n\n"
	detail := bodyStyle.Render(c.Detail)

	url := ""
	if c.URL != "" {
		url = "\n\n" + accentStyle.Render("→ "+c.URL)
	}

	footer := footerStyle.Render("\n\nesc go back")
	return lipgloss.NewStyle().Padding(1, 2).Render(header + category + desc + detail + url + footer)
}

var lightningFrames = []string{
	// dim — dots
	"    ....\n" +
		"   ..\n" +
		"  ..\n" +
		" ......\n" +
		"   ..\n" +
		"    ..\n" +
		"     .",
	// medium — plus
	"    ++++\n" +
		"   ++\n" +
		"  ++\n" +
		" ++++++\n" +
		"   ++\n" +
		"    ++\n" +
		"     +",
	// bright — stars
	"    ****\n" +
		"   **\n" +
		"  **\n" +
		" ******\n" +
		"   **\n" +
		"    **\n" +
		"     *",
	// flash — hash
	"    ####\n" +
		"   ##\n" +
		"  ##\n" +
		" ######\n" +
		"   ##\n" +
		"    ##\n" +
		"     #",
}

func viewReflections(m Model) string {
	header := boldStyle.Render("Reflections") + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	var body strings.Builder
	for i, r := range allReflections {
		if i == m.reflectionIndex {
			body.WriteString(selectStyle.Render("  ✦ "+r.Title) + "\n")
		} else {
			body.WriteString(subtleStyle.Render("    "+r.Title) + "\n")
		}
	}

	footer := footerStyle.Render("\n↑/↓ navigate • enter open • esc back")

	bolt := accentStyle.Render(lightningFrames[m.frameIndex%len(lightningFrames)])

	left := lipgloss.NewStyle().Width(46).Render(header + body.String() + footer)
	right := lipgloss.NewStyle().Padding(1, 8).Render(bolt)

	return lipgloss.NewStyle().Padding(1, 2).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
	)
}

func viewReflectionDetail(m Model) string {
	r := allReflections[m.reflectionIndex]

	header := boldStyle.Render(r.Title) + "\n" +
		accentStyle.Render(strings.Repeat("─", 30)) + "\n\n"

	detail := bodyStyle.Render(r.Detail)

	url := ""
	if r.URL != "" {
		url = "\n\n" + accentStyle.Render("→ "+r.URL)
	}

	footer := footerStyle.Render("\n\nesc go back")
	return lipgloss.NewStyle().Padding(1, 2).Render(header + detail + url + footer)
}


var contactStarFrames = []string{
	"    .    \n" +
		"    |    \n" +
		"  . | .  \n" +
		". --+-- .\n" +
		"  ' | '  \n" +
		"    |    \n" +
		"    '    ",
	"    *    \n" +
		"    |    \n" +
		"  * | *  \n" +
		"* --+-- *\n" +
		"  * | *  \n" +
		"    |    \n" +
		"    *    ",
	"    ✦    \n" +
		"    |    \n" +
		"  ✦ | ✦  \n" +
		"✦ --✦-- ✦\n" +
		"  ✦ | ✦  \n" +
		"    |    \n" +
		"    ✦    ",
	"    *    \n" +
		"    |    \n" +
		"  + | +  \n" +
		"+ --*-- +\n" +
		"  + | +  \n" +
		"    |    \n" +
		"    *    ",
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

	star := accentStyle.Render(contactStarFrames[m.frameIndex%len(contactStarFrames)])

	left := lipgloss.NewStyle().Width(36).Render(header + links + footer)
	right := lipgloss.NewStyle().Padding(1, 8).Render(star)

	return lipgloss.NewStyle().Padding(1, 2).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
	)
}

// project list

