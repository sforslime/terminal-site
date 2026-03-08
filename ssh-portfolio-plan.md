# SSH Terminal Portfolio — Implementation Plan

## The Idea

Visitors type `ssh mori.dev` in their terminal and land in a fully interactive, beautifully rendered terminal UI — ASCII art portrait, bio, navigable sections — no browser needed.

---

## Architecture Overview

```
┌──────────────────────────────────────────────┐
│  Visitor's Terminal                           │
│  $ ssh mori.dev                              │
└──────────────┬───────────────────────────────┘
               │ SSH connection
               ▼
┌──────────────────────────────────────────────┐
│  Your VPS (e.g. DigitalOcean, Fly.io)        │
│                                              │
│  ┌────────────────────────────────────────┐  │
│  │  Wish (SSH Server)                     │  │
│  │  Listens on port 22 (or 2222)          │  │
│  │  No auth required — public portfolio   │  │
│  └──────────────┬─────────────────────────┘  │
│                 │                             │
│  ┌──────────────▼─────────────────────────┐  │
│  │  Bubble Tea TUI App                    │  │
│  │                                        │  │
│  │  ┌─────────┐ ┌──────────┐ ┌────────┐  │  │
│  │  │ Home    │ │Creations │ │Contact │  │  │
│  │  │(ASCII + │ │(project  │ │(links) │  │  │
│  │  │ bio)    │ │ list)    │ │        │  │  │
│  │  └─────────┘ └──────────┘ └────────┘  │  │
│  │                                        │  │
│  │  Lip Gloss (styling/layout)            │  │
│  └────────────────────────────────────────┘  │
│                                              │
└──────────────────────────────────────────────┘
```

---

## Tech Stack

| Layer            | Tool                          | Why                                      |
|------------------|-------------------------------|------------------------------------------|
| Language         | **Go**                        | Single binary, fast, Charm ecosystem     |
| SSH Server       | **Wish** (charmbracelet/wish) | Turns any Go app into an SSH server      |
| TUI Framework    | **Bubble Tea**                | Elm-architecture TUI, handles input/rendering |
| Styling          | **Lip Gloss**                 | Terminal CSS — colors, borders, layout    |
| ASCII Art        | **ascii-image-converter**     | Convert your photo to ASCII at build time|
| Hosting          | **Fly.io** or **DigitalOcean**| Cheap VPS with a public IP               |
| Domain           | Any registrar                 | Point an A record to your VPS            |

---

## Phase 1: Local Bubble Tea App (Days 1–2)

Build the TUI app first, test it locally in your own terminal.

### Task 1.1 — Project Setup

```bash
mkdir ssh-portfolio && cd ssh-portfolio
go mod init github.com/yourname/ssh-portfolio
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/wish
go get github.com/charmbracelet/wish/bubbletea
```

**File structure:**

```
ssh-portfolio/
├── main.go              # SSH server entry point
├── tui/
│   ├── model.go         # Bubble Tea model (state + logic)
│   ├── views.go         # Render functions for each screen
│   ├── styles.go        # Lip Gloss style definitions
│   └── constants.go     # ASCII art, bio text, project data
├── ascii/
│   └── portrait.txt     # Your pre-generated ASCII portrait
├── Dockerfile
├── fly.toml             # Fly.io config (if using Fly)
└── go.mod
```

### Task 1.2 — Define the Bubble Tea Model

The core state machine. Bubble Tea uses the Elm architecture: Model → Update → View.

```go
// tui/model.go
package tui

type screen int

const (
    screenHome screen = iota
    screenCreations
    screenReflections
    screenContacts
)

type model struct {
    currentScreen screen
    navIndex      int      // which nav item is highlighted
    navItems      []string // ["Creations", "Reflections", "Contacts"]
    width         int      // terminal width (from WindowSizeMsg)
    height        int      // terminal height
}
```

Key messages to handle:
- `tea.KeyMsg` — arrow keys to navigate, enter to select, `q` to quit
- `tea.WindowSizeMsg` — adapt layout to visitor's terminal size

### Task 1.3 — Build the Home View

This is the main screen from the screenshot. Use Lip Gloss for layout:

```
┌─────────────────────────┬──────────────────────────┐
│                         │  * ASCII logo "MORI"     │
│   ASCII portrait        │                          │
│   (loaded from .txt)    │  Bio text                │
│                         │  (dim text for older bio) │
│                         │                          │
│                         │  [Creations] +Reflections │
│                         │   Contacts               │
├─────────────────────────┴──────────────────────────┤
│ [← → to select · enter to open · q to quit]       │
└────────────────────────────────────────────────────┘
```

Lip Gloss techniques:
- `lipgloss.JoinHorizontal()` — place portrait and bio side by side
- `lipgloss.NewStyle().Foreground(lipgloss.Color("#4ecdc4"))` — accent color
- `lipgloss.NewStyle().Faint(true)` — dim text for secondary bio
- `lipgloss.Place()` — center the whole layout in the terminal

### Task 1.4 — Generate Your ASCII Portrait

Use the `ascii-image-converter` CLI tool:

```bash
# Install
go install github.com/TheZoraworIO/ascii-image-converter@latest

# Generate — experiment with these flags
ascii-image-converter your-photo.jpg \
  --width 60 \
  --dither \
  --color-bg \
  --map " .:-=+*#%@"
```

Save the output to `ascii/portrait.txt` and load it at startup.
Tip: Use a high-contrast photo with a clean background for best results.

### Task 1.5 — Build Sub-Screens

Each nav section gets its own view function:

**Creations** — A scrollable list of projects:
```
┌────────────────────────────────────┐
│  > project-name                    │
│    A short description             │
│    github.com/you/project          │
│                                    │
│    another-project                 │
│    Description here                │
│    link.com                        │
└────────────────────────────────────┘
```

Use `charmbracelet/bubbles/list` for a pre-built interactive list component.

**Reflections** — Blog post titles or thoughts (link out or display inline).

**Contacts** — Simple list of links (GitHub, Twitter, email).

### Task 1.6 — Test Locally

```bash
go run main.go
# This runs the Bubble Tea app directly in your terminal
# Verify: navigation works, layout adapts to resize, q quits
```

---

## Phase 2: Wrap in SSH Server (Day 3)

### Task 2.1 — Set Up Wish

Wish makes your Bubble Tea app servable over SSH with ~30 lines:

```go
// main.go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/wish"
    bm "github.com/charmbracelet/wish/bubbletea"
    "github.com/charmbracelet/ssh"
    "github.com/yourname/ssh-portfolio/tui"
)

const host = "0.0.0.0"
const port = 2222

func main() {
    s, err := wish.NewServer(
        wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
        wish.WithHostKeyPath(".ssh/term_info_ed25519"),
        wish.WithMiddleware(
            bm.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
                pty, _, _ := s.Pty()
                m := tui.NewModel(pty.Window.Width, pty.Window.Height)
                return m, []tea.ProgramOption{tea.WithAltScreen()}
            }),
        ),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not start server: %s\n", err)
        os.Exit(1)
    }

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err = s.ListenAndServe(); err != nil {
            fmt.Fprintf(os.Stderr, "Server error: %s\n", err)
        }
    }()

    fmt.Printf("SSH server running on %s:%d\n", host, port)
    <-done

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    s.Shutdown(ctx)
}
```

### Task 2.2 — Test Over SSH Locally

```bash
go run main.go &
ssh localhost -p 2222
# You should see your full portfolio TUI!
```

---

## Phase 3: Deploy & Go Live (Day 4)

### Option A: Fly.io (Recommended — easy, cheap, global)

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o portfolio .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/portfolio .
COPY ascii/ ./ascii/
EXPOSE 2222
CMD ["./portfolio"]
```

```toml
# fly.toml
app = "mori-portfolio"

[build]
  dockerfile = "Dockerfile"

[[services]]
  internal_port = 2222
  protocol = "tcp"

  [[services.ports]]
    port = 22          # So visitors can just `ssh mori.dev`
```

```bash
fly launch
fly deploy
fly ips allocate-v4   # Get a static IP
```

### Option B: DigitalOcean / Any VPS

```bash
# On your VPS
scp ./portfolio yourserver:/opt/portfolio/
ssh yourserver

# Create a systemd service
sudo cat > /etc/systemd/system/portfolio.service << EOF
[Unit]
Description=SSH Portfolio
After=network.target

[Service]
ExecStart=/opt/portfolio/portfolio
Restart=always
User=portfolio
WorkingDirectory=/opt/portfolio

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable portfolio
sudo systemctl start portfolio
```

If you want it on port 22 (so visitors don't need `-p 2222`), either:
- Move your real SSH to port 2222 and give the portfolio port 22
- Use iptables to route port 22 → 2222 for non-authenticated connections

### Task 3.1 — DNS Setup

Point your domain's **A record** to your server's IP:

```
mori.dev  →  A  →  143.xxx.xxx.xxx
```

Now `ssh mori.dev` just works.

---

## Phase 4: Polish & Extras (Ongoing)

### Nice-to-haves

- **Logging/analytics** — Count unique SSH connections (log `s.RemoteAddr()`)
- **Adaptive layout** — Already built in via `WindowSizeMsg`, but test on small terminals (80x24)
- **Loading animation** — Show a brief typewriter effect or progress bar on connect
- **Color detection** — Check `s.Pty().Term` and degrade gracefully for terminals without color
- **Rate limiting** — Prevent abuse with `wish/ratelimiter` middleware
- **Idle timeout** — Disconnect after 5 min of inactivity

### Inspiration & References

| Project                | What to learn from it           |
|------------------------|---------------------------------|
| `ssh chat.shazow.net`  | Multi-user SSH app              |
| `ssh git.charm.sh`     | Charm's own SSH TUI             |
| `ssh caarlos0.dev`     | Carlos Becker's SSH portfolio   |
| Bubble Tea examples    | github.com/charmbracelet/bubbletea/examples |
| Lip Gloss examples     | github.com/charmbracelet/lipgloss |

---

## Quick Start Checklist

- [ ] Install Go 1.22+
- [ ] `go mod init` and pull Charm dependencies
- [ ] Generate ASCII portrait from your photo
- [ ] Build model.go with screen states and navigation
- [ ] Build home view with Lip Gloss (portrait + bio side-by-side)
- [ ] Build sub-screen views (Creations, Reflections, Contacts)
- [ ] Wrap in Wish SSH server
- [ ] Test locally with `ssh localhost -p 2222`
- [ ] Dockerize and deploy to Fly.io or VPS
- [ ] Point DNS A record to server IP
- [ ] `ssh mori.dev` — you're live
