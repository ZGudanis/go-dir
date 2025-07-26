package godir

import tea "github.com/charmbracelet/bubbletea"

type DirectoryModel struct {
}

func (m DirectoryModel) Init() tea.Cmd {
	return nil
}

func (m DirectoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m DirectoryModel) View() string {
	return ""
}
