package ttt

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

/*
TicTacToe
*/

type Model struct {
	gs       gameState
	Quitting bool // User quit?
}

type gameState struct {
	Cells []string // 1d repro of 2d array
	Rows  int      // Slice helper
	Cols  int      // Slice helper
}

func NewModel() Model {
	// TODO: Get board size as args
	return Model{
		gs: gameState{
			Cells: make([]string, 3*3),
			Rows:  3,
			Cols:  3,
		},
		Quitting: false,
	}
}

type TTT interface {
	// SetAt sets a cell to a value (x/o/empty)
	SetAt(row int, col int, value string) error
	// GetAt returns the value of a cell
	GetAt(row int, col int) (string, error)
}

func (m Model) SetAt(row int, col int, value string) error {
	if row > m.gs.Rows || col > m.gs.Cols {
		return fmt.Errorf("out of bounds")
	}

	m.gs.Cells[(col*row)-1] = value

	return nil
}

func (m Model) GetAt(row int, col int) (string, error) {
	if row > m.gs.Rows || col > m.gs.Cols {
		return "", fmt.Errorf("out of bounds")
	}

	return m.gs.Cells[(row*col)-1], nil
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	/*
		// Hand off the message and model to the appropriate update function for the
		// appropriate view based on the current state.
		if !m.Chosen {
			return updateChoices(msg, m)
		}
		return updateChosen(msg, m) */

	return m, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m Model) View() string {
	return fmt.Sprintf("This is the view!\n%+v", m)
}
