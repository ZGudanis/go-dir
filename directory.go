package godir

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type DirectoryModel struct {
	cwd         string
	cursor      int
	entries     []os.DirEntry
	showPreview bool

	search   string
	isSearch bool
}

func InitModel(cwd string) DirectoryModel {
	entries, err := os.ReadDir(cwd)
	if err != nil {
		panic(err)
	}
	return DirectoryModel{
		entries: entries,
		cwd:     cwd,
	}
}

func (m DirectoryModel) Init() tea.Cmd {
	return nil
}

func (m DirectoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.keyMsgHandler(msg)
	}

	return m, nil
}

func (m DirectoryModel) View() string {
	sbuilder := &strings.Builder{}

	sbuilder.WriteString(RenderList(m))

	return sbuilder.String()
}

type MoveCmd int

const (
	Forward MoveCmd = iota
	Backward
	Into
	Outof
)

// var moveCmdName = map[MoveCmd]string{
// 	Forward:  "forward",
// 	Backward: "backward",
// 	Into:     "into",
// 	Outof:    "outof",
// }

func updatePath(model DirectoryModel, move MoveCmd) DirectoryModel {
	var path string
	switch move {
	case Into:
		if len(model.cwd) == 1 {
			path = model.cwd + model.entries[model.cursor].Name()
		} else {
			path = model.cwd + "/" + model.entries[model.cursor].Name()
		}
	case Outof:
		index := strings.LastIndex(model.cwd, "/")
		if index != 0 {
			path = model.cwd[0:index]
		} else {
			path = "/"
		}
	default:
		return model
	}

	dir, err := os.ReadDir(path)
	if err != nil {
		return model
	}
	return DirectoryModel{
		entries:     dir,
		cursor:      0,
		cwd:         path,
		showPreview: model.showPreview,
	}
}

func (m DirectoryModel) keyMsgHandler(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "ctrl+f": // Search{}
		m.isSearch = false
		searchField.Blur()
		return m, nil

	case "ctrl+v":
		if m.entries[m.cursor].IsDir() {
			return m, nil
		}
		return m, tea.ExecProcess(exec.Command("vim", m.cwd+"/"+m.entries[m.cursor].Name()), nil)
	case "ctrl+n":
		if m.entries[m.cursor].IsDir() {
			return m, nil
		}
		return m, tea.ExecProcess(exec.Command("nano", m.cwd+"/"+m.entries[m.cursor].Name()), nil)

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			return m, nil
		}

	case "down", "j":
		if m.cursor < len(m.entries)-1 {
			m.cursor++
			return m, nil
		}

	case "ctrl+x": // Search{}
		searchField.Reset()
		return m, nil

	case "pgdown":
		if len(m.entries) > listHeight && m.cursor < len(m.entries)-listHeight {
			m.cursor += listHeight
		} else {
			m.cursor = len(m.entries) - 1
		}

		return m, nil

	case "pgup":
		if m.cursor-listHeight > 0 {
			m.cursor -= listHeight
		} else {
			m.cursor = 0
		}
		return m, nil
	}

	var cmd tea.Cmd
	if searchField.Focused() {
		searchField, cmd = searchField.Update(key)
		return m, cmd
	}

	switch key.String() {
	case "r": // Refresh
		return m, nil
	case "p":
		m.showPreview = !m.showPreview
		return m, nil
	case "q":
		return m, tea.Quit
	case "f": // Search{}
		m.isSearch = true
		return m, searchField.Focus()

	case "enter":
		if m.isSearch {
			m.isSearch = !m.isSearch
			return m, nil
		}
		return m, tea.ExecProcess(exec.Command("code", m.cwd), nil)

	case "ctrl+e":
		if len(m.cwd) == 1 {
			return m, tea.ExecProcess(exec.Command("code", m.cwd+m.entries[m.cursor].Name()), nil)
		}
		return m, tea.ExecProcess(exec.Command("code", m.cwd+"/"+m.entries[m.cursor].Name()), nil)

	case "right":
		if len(m.entries) == 0 {
			return m, nil
		}
		return updatePath(m, Into), nil

	case "left":
		return updatePath(m, Outof), nil

	}

	return m, nil
}
