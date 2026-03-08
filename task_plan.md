# SSH Terminal Portfolio — Task Plan

## Project Goal
Build an interactive terminal portfolio accessible via `ssh <your-domain>`. Visitors connect and land in a full Bubble Tea TUI — ASCII portrait, bio, navigable sections — no browser required.

## Status: IN PROGRESS
Last updated: 2026-03-08

---

## Phase 1: Local Bubble Tea App
**Goal:** Build and test the full TUI locally before adding SSH.

### Tasks

- [ ] **1.1 Project Setup**
  - `go mod init github.com/sforslime/ssh-portfolio`
  - Install deps: bubbletea, lipgloss, wish, wish/bubbletea, bubbles
  - Create file structure: `main.go`, `tui/`, `ascii/`, `Dockerfile`, `fly.toml`
  - Status: pending

- [ ] **1.2 Define Bubble Tea Model**
  - `tui/model.go` — screen states (Home, Creations, Reflections, Contacts)
  - Navigation state: `currentScreen`, `navIndex`, `navItems`, `width`, `height`
  - Handle `tea.KeyMsg` (arrows, enter, q) and `tea.WindowSizeMsg`
  - Status: pending

- [ ] **1.3 Build Home View**
  - `tui/views.go` — two-column layout: ASCII portrait left, bio + nav right
  - Lip Gloss: `JoinHorizontal`, `Place`, accent color `#4ecdc4`, faint secondary bio
  - Footer: key hints bar
  - Status: pending

- [ ] **1.4 Generate ASCII Portrait**
  - Install `ascii-image-converter`
  - Run with `--width 60 --dither --color-bg` flags on a high-contrast photo
  - Save output to `ascii/portrait.txt`
  - Status: pending

- [ ] **1.5 Build Sub-Screen Views**
  - Creations: scrollable project list using `bubbles/list`
  - Reflections: blog titles or inline thoughts
  - Contacts: GitHub, Twitter, email links
  - Status: pending

- [ ] **1.6 Test Locally**
  - `go run main.go` — verify navigation, layout resize, q to quit
  - Status: pending

---

## Phase 2: SSH Server Wrapper
**Goal:** Expose the TUI over SSH using Wish.

### Tasks

- [ ] **2.1 Set Up Wish Server**
  - `main.go` — Wish server on `0.0.0.0:2222`
  - Host key path: `.ssh/term_info_ed25519`
  - Middleware: `bm.Middleware` passing PTY size to `tui.NewModel`
  - `tea.WithAltScreen()` program option
  - Graceful shutdown on SIGINT/SIGTERM
  - Status: pending

- [ ] **2.2 Test Over SSH Locally**
  - `go run main.go &` then `ssh localhost -p 2222`
  - Verify full TUI renders correctly over SSH
  - Status: pending

---

## Phase 3: Deploy & Go Live
**Goal:** Ship to a public server and wire up DNS.

### Tasks

- [ ] **3.1 Dockerize**
  - Multi-stage Dockerfile: Go builder → alpine runtime
  - Copy binary + `ascii/` dir, expose port 2222
  - Status: pending

- [ ] **3.2 Deploy to Fly.io** (recommended)
  - `fly.toml`: expose internal 2222 → external 22
  - `fly launch` → `fly deploy` → `fly ips allocate-v4`
  - Status: pending

- [ ] **3.3 DNS Setup**
  - A record: `<your-domain> → <fly static IP>`
  - Verify `ssh <your-domain>` connects
  - Status: pending

---

## Phase 4: Polish & Extras
**Goal:** Production hardening and nice-to-haves.

### Tasks

- [ ] **4.1 Rate Limiting** — `wish/ratelimiter` middleware
- [ ] **4.2 Idle Timeout** — disconnect after 5 min inactivity
- [ ] **4.3 Analytics/Logging** — log `s.RemoteAddr()` per connection
- [ ] **4.4 Loading Animation** — typewriter effect or progress bar on connect
- [ ] **4.5 Color Detection** — check `s.Pty().Term`, degrade gracefully
- [ ] **4.6 Small Terminal Testing** — verify layout at 80x24

---

## Key Decisions

| Decision | Choice | Reason |
|---|---|---|
| Language | Go | Single binary, Charm ecosystem |
| SSH server | Wish | Minimal setup, Bubble Tea integration |
| TUI framework | Bubble Tea | Elm-arch, handles SSH rendering |
| Styling | Lip Gloss | Terminal CSS |
| Hosting | Fly.io | Easy, cheap, global TCP routing |
| Port | 22 (via Fly routing) | No `-p` flag needed for visitors |
