# Findings & Research

## Charm Ecosystem

### Wish
- `charmbracelet/wish` — turns any Go app into an SSH server in ~30 lines
- `wish/bubbletea` middleware (`bm.Middleware`) bridges SSH PTY to Bubble Tea program
- `wish/ratelimiter` — built-in middleware for abuse prevention
- Host key auto-generated on first run if path doesn't exist
- No auth required by default — perfect for public portfolio

### Bubble Tea
- Elm architecture: `Model`, `Update(msg) (Model, Cmd)`, `View() string`
- Key messages: `tea.KeyMsg`, `tea.WindowSizeMsg`
- `tea.WithAltScreen()` — uses alternate terminal buffer (clean, no scroll history bleed)
- `charmbracelet/bubbles/list` — pre-built scrollable list component for Creations screen

### Lip Gloss
- `lipgloss.JoinHorizontal(lipgloss.Top, left, right)` — side-by-side columns
- `lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)` — center in terminal
- `lipgloss.NewStyle().Foreground(lipgloss.Color("#4ecdc4"))` — accent color
- `lipgloss.NewStyle().Faint(true)` — dim secondary text
- `lipgloss.NewStyle().Border(lipgloss.RoundedBorder())` — rounded box borders

## ASCII Portrait Generation
- Tool: `github.com/TheZoraworIO/ascii-image-converter`
- Best flags: `--width 60 --dither --color-bg --map " .:-=+*#%@"`
- Photo tips: high contrast, clean/simple background, face centered
- Output saved to `ascii/portrait.txt`, loaded at startup

## Fly.io Deployment Notes
- TCP routing: map external port 22 → internal container port 2222
- `fly ips allocate-v4` — get static IPv4 for DNS A record
- `fly.toml` `[[services]]` with `protocol = "tcp"` for raw TCP (not HTTP)
- Free tier has limits; SSH portfolio traffic is minimal

## Inspiration Projects
- `ssh caarlos0.dev` — Carlos Becker's SSH portfolio (direct reference)
- `ssh chat.shazow.net` — multi-user SSH chat
- `ssh git.charm.sh` — Charm's own SSH TUI

## Potential Gotchas
- Terminal color support varies — check `TERM` env or `s.Pty().Term`
- Small terminals (80x24) may need layout adjustments — test explicitly
- Port 22 on a VPS conflicts with your own SSH — Fly.io sidesteps this cleanly
- ASCII art with `--color-bg` uses ANSI color codes; width measurement may be off — use `lipgloss.Width()` for accurate measurement
