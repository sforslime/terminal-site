package tui

import (
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// flashStar randomly swaps ~8% of characters with similar alternatives to create a shimmer effect.
func flashStar(base string) string {
	swaps := map[rune][]rune{
		'.': {',', '\'', '`'},
		',': {'.', '\''},
		'z': {'s', 'x', 'Z'},
		'Z': {'X', 'Y', 'z'},
		'X': {'Y', 'Z', 'x'},
		'x': {'z', 'X'},
		'[': {'{', '('},
		']': {'}', ')'},
		'{': {'[', '('},
		'}': {']', ')'},
		'v': {'u', 'w'},
		'n': {'m', 'r'},
	}
	runes := []rune(base)
	for i, r := range runes {
		if alts, ok := swaps[r]; ok && rand.Float32() < 0.08 {
			runes[i] = alts[rand.Intn(len(alts))]
		}
	}
	return string(runes)
}

// loadFrames reads all .txt files from ascii/frames/ sorted by filename.
// Falls back to the static portrait if no frames exist.
func loadFrames() []string {
	entries, err := os.ReadDir("ascii/frames")
	if err != nil || len(entries) == 0 {
		// fall back to static portrait
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

type Model struct {
	screen     screen
	navIndex   int
	width      int
	height     int
	list       list.Model
	snow       []snowflake
	frames     []string
	frameIndex int
	starBase   string
	star       string
}

func NewModel(width, height int) Model {
	starBase := ""
	if data, err := os.ReadFile("ascii/star.txt"); err == nil {
		starBase = string(data)
	}
	return Model{
		screen:   screenHome,
		width:    width,
		height:   height,
		list:     newProjectList(),
		snow:     initSnow(),
		frames:   loadFrames(),
		starBase: starBase,
		star:     flashStar(starBase),
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
		m.star = flashStar(m.starBase)
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
