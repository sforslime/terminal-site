package tui

import (
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenHome screen = iota
	screenCreations
	screenReflections
	screenContacts
	screenCreationDetail
	screenReflectionDetail
)

var navItems = []string{"Creations", "Reflections", "Contacts"}

// Creation data

type Creation struct {
	Category string
	Title    string
	Desc     string
	Detail   string
	URL      string
}

var allCreations = []Creation{
	{
		Category: "Personal Projects",
		Title:    "ssh-portfolio",
		Desc:     "interactive terminal portfolio over SSH",
		Detail:   "Built with Go, Bubble Tea, Lip Gloss, and Wish.\nVisitors SSH in and get a full TUI — no browser needed.\nAnimated ASCII art, navigable sections, snow effect.",
		URL:      "github.com/sforslime/ssh-portfolio",
	},
	{
		Category: "Personal Projects",
		Title:    "web scraper",
		Desc:     "data scraper built with Selenium and Python",
		Detail:   "Scrapes structured data from dynamic web pages.\nUses Selenium WebDriver with headless Chrome.\nOutputs clean JSON and CSV.",
		URL:      "github.com/sforslime/web-scraper",
	},
	{
		Category: "AI Research",
		Title:    "project three",
		Desc:     "add your project here",
		Detail:   "Add more details about this project here.",
		URL:      "",
	},
}

// Reflection data

type Reflection struct {
	Title  string
	Detail string
	URL    string
}

var allReflections = []Reflection{
	{
		Title:  "building in public is underrated",
		Detail: "Sharing your process — the messy drafts, the dead ends, the small wins —\nbuilds trust faster than any polished launch post ever could.\nPeople connect with the journey, not just the destination.",
		URL:    "",
	},
	{
		Title:  "terminals are timeless",
		Detail: "GUIs come and go, but the terminal endures.\nThere's something honest about a blank prompt —\nno dark patterns, no infinite scroll, just you and the machine.",
		URL:    "",
	},
	{
		Title:  "simplicity beats complexity, always",
		Detail: "Every abstraction has a cost. Every added layer is a future burden.\nThe best systems are the ones you can hold entirely in your head.\nSimplicity isn't laziness — it's discipline.",
		URL:    "",
	},
}

// Snow

const snowW = 48
const snowH = 10

type snowflake struct {
	x, y int
	ch   rune
}

var snowChars = []rune{'*', '.', '\'', ',', '+', '·'}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
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

// Frames

func loadFrames() []string {
	entries, err := os.ReadDir("ascii/frames")
	if err != nil || len(entries) == 0 {
		data, err := os.ReadFile("ascii/portrait.txt")
		if err != nil {
			return []string{"[ no portrait ]"}
		}
		return []string{string(data)}
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".txt") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	frames := make([]string, 0, len(names))
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join("ascii/frames", name))
		if err == nil {
			frames = append(frames, string(data))
		}
	}
	if len(frames) == 0 {
		return []string{"[ no frames ]"}
	}
	return frames
}

// Model

type Model struct {
	screen          screen
	navIndex        int
	width           int
	height          int
	snow            []snowflake
	frames          []string
	frameIndex      int
	creationIndex   int
	reflectionIndex int
}

func NewModel(width, height int) Model {
	return Model{
		screen: screenHome,
		width:  width,
		height: height,
		snow:   initSnow(),
		frames: loadFrames(),
	}
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.snow = tickSnow(m.snow)
		m.frameIndex = (m.frameIndex + 1) % len(m.frames)
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
				if m.creationIndex > 0 {
					m.creationIndex--
				}
			} else if m.screen == screenReflections {
				if m.reflectionIndex > 0 {
					m.reflectionIndex--
				}
			}

		case "down", "j":
			if m.screen == screenCreations {
				if m.creationIndex < len(allCreations)-1 {
					m.creationIndex++
				}
			} else if m.screen == screenReflections {
				if m.reflectionIndex < len(allReflections)-1 {
					m.reflectionIndex++
				}
			}

		case "enter":
			if m.screen == screenHome {
				m.screen = screen(m.navIndex + 1)
			} else if m.screen == screenCreations {
				m.screen = screenCreationDetail
			} else if m.screen == screenReflections {
				m.screen = screenReflectionDetail
			}

		case "esc", "backspace":
			if m.screen == screenCreationDetail {
				m.screen = screenCreations
			} else if m.screen == screenReflectionDetail {
				m.screen = screenReflections
			} else {
				m.screen = screenHome
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case screenCreations:
		return viewCreations(m)
	case screenCreationDetail:
		return viewCreationDetail(m)
	case screenReflections:
		return viewReflections(m)
	case screenReflectionDetail:
		return viewReflectionDetail(m)
	case screenContacts:
		return viewContacts(m)
	default:
		return viewHome(m)
	}
}
