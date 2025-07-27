package main

import (
	"fmt"
	"os"
	"path/filepath"

	godir "github.com/ZGudanis/go-dir"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(godir.InitModel(filepath.Dir(ex)))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v", err)
		os.Exit(1)
	}
}

type Model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initModel() Model {
	return Model{
		choices:  []string{"Hello", "This is", "Me again"},
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			{
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "down", "j":
			{
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			}
		case "enter", " ":
			{
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := "What should i do\n\n"
	for i, choices := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choices)
	}

	return s + "\nPress q to exit\n"
}
