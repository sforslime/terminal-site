# ssh-portfolio

A portfolio you connect to with `ssh`, not a browser. Visitors drop into a Bubble Tea TUI with an animated ASCII self-portrait, navigable Creations and Reflections sections, and contact links — all served as a single Go binary via Wish.


https://github.com/user-attachments/assets/89ec764b-4f00-4b45-85f5-39c0b6032708



No ssh domain yet though😔


## Stack

- **Go** — single binary, fast startup
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — Elm-architecture TUI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** — terminal styling
- **[Wish](https://github.com/charmbracelet/wish)** — SSH server wrapping Bubble Tea

## Features

- Two-column home screen: animated ASCII portrait (left), bio + navigation (right)
- Snow animation drifting through the AYO! logo
- **Creations** — categorized project list with detail view on enter
- **Reflections** — navigable thoughts with expandable detail view
- **Contacts** — links with animated ASCII star
- Animated video frames cycling as the portrait (braille + dither mode)
- Graceful shutdown on SIGINT/SIGTERM

## Navigation

| Key | Action |
|-----|--------|
| `←` `→` | Move between nav items on home |
| `↑` `↓` | Navigate lists (Creations, Reflections) |
| `enter` | Open selected item |
| `esc` / `backspace` | Go back |
| `q` / `ctrl+c` | Quit |

## Project Structure

```
.
├── main.go              # Wish SSH server
├── tui/
│   ├── model.go         # Bubble Tea model, state, data
│   └── views.go         # All view rendering
├── ascii/
│   ├── portrait.txt     # Fallback static ASCII portrait
│   └── frames/          # Animated portrait frames (frame_0001.txt …)
└── .ssh/
    └── term_info_ed25519  # SSH host key
```

## Running Locally

```bash
go run main.go
```

Then in another terminal:

```bash
ssh localhost -p 2222
```

## Generating Portrait Frames

Extract frames from a video and convert to ASCII (braille + dither):

```bash
# Extract frames at 10fps, scaled to 120px wide
ffmpeg -i ascii/video.mp4 -vf "fps=10,scale=120:-1" ascii/frames/frame_%04d.png

# Convert each frame to ASCII
for f in ascii/frames/*.png; do
  ascii-image-converter "$f" -W 60 -b --dither -o "${f%.png}.txt"
done
```

## Customizing Content

Edit `tui/model.go`:

- **`allCreations`** — add/edit projects (Category, Title, Desc, Detail, URL)
- **`allReflections`** — add/edit thoughts (Title, Detail, URL)
- Contact links are in `viewContacts` inside `tui/views.go`

## Deployment (Fly.io)

```bash
fly launch
fly deploy
fly ips allocate-v4
```

Map an A record from your domain to the allocated IP. The `fly.toml` should route external port 22 → internal 2222.

## License

MIT
