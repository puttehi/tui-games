package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/puttehi/tui-games/internal/ttt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

const (
	host = "localhost"
	port = 6969
)

// app contains a wish server and the list of running programs.
type app struct {
	*ssh.Server
	progs []*tea.Program // Client programs I think?
}

// send dispatches a message to all running programs.
func (a *app) send(msg tea.Msg) {
	for _, p := range a.progs {
		go p.Send(msg)
	}
}

func newApp() *app {
	a := new(app)
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/host_key_ed25519"),
		wish.WithMiddleware(
			bm.MiddlewareWithProgramHandler(a.ProgramHandler, termenv.ANSI256),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	a.Server = s
	return a
}

func (a *app) Start() {
	var err error
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		if err = a.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := a.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

func (a *app) ProgramHandler(s ssh.Session) *tea.Program {
	if _, _, active := s.Pty(); !active {
		wish.Fatalln(s, "terminal is not active")
	}

	m := model{
		app: a,
		ttt: ttt.NewModel(), // How the heck is View() coming from here (ttt.Model.View()) and not from model.View() ?!
	}

	fmt.Println("PROGRAM HANDLER WTF")
	fmt.Println(m.View())

	p := tea.NewProgram(m, bm.MakeOptions(s)...)
	a.progs = append(a.progs, p)

	return p
}

func main() {
	app := newApp()
	app.Start()
}

type model struct {
	*app
	ttt ttt.Model
}

func (m model) Init() tea.Cmd {
	return m.ttt.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// we might need to do server stuff and send events to progs here with app.send but lets see
	return m.ttt.Update(msg)
}

func (m model) View() string {
	fmt.Println("VIEW()")
	return "WTF"
	// return fmt.Sprintf("Hello from server, clients: %d\n%s", len(m.app.progs), m.ttt.View())
}
