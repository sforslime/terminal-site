package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	cssh "github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/sforslime/ssh-portfolio/tui"
)

const (
	host    = "0.0.0.0"
	port    = 2222
	keyPath = ".ssh/term_info_ed25519"
)

func main() {
	srv, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create server: %v\n", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("SSH server listening on %s:%d\n", host, port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	<-done

	fmt.Println("\nShutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}
}

func teaHandler(s cssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	w := pty.Window.Width
	h := pty.Window.Height
	if w == 0 {
		w = 120
	}
	if h == 0 {
		h = 40
	}
	m := tui.NewModel(w, h)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
