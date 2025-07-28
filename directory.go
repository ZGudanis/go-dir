package godir

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type DirectoryModel struct {
	cursor  int
	entries []os.DirEntry
	cwd     string
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
		entries: dir,
		cursor:  0,
		cwd:     path,
	}
}

func (m DirectoryModel) keyMsgHandler(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "enter":
		return m, tea.ExecProcess(exec.Command("code", m.cwd), nil)
	case "ctrl+e":
		return m, tea.ExecProcess(exec.Command("code", m.cwd+"/"+m.entries[m.cursor].Name()), nil)
	case "ctrl+v":
		{
			if m.entries[m.cursor].IsDir() {
				return m, nil
			}
			return m, tea.ExecProcess(exec.Command("vim", m.cwd+"/"+m.entries[m.cursor].Name()), nil)
		}
	case "ctrl+n":
		{
			if m.entries[m.cursor].IsDir() {
				return m, nil
			}
			return m, tea.ExecProcess(exec.Command("nano", m.cwd+"/"+m.entries[m.cursor].Name()), nil)
		}

	case "up", "k":
		{
			if m.cursor > 0 {
				m.cursor--
			}
		}

	case "down", "j":
		{
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		}
	case "right":
		{
			if len(m.entries) == 0 {
				return m, nil
			}
			return updatePath(m, Into), nil
		}
	case "left":
		return updatePath(m, Outof), nil

	case "pgdown":
		{
			if len(m.entries) > listHeight && m.cursor < len(m.entries)-listHeight {

				m.cursor += listHeight
				return m, nil
			}
			m.cursor = len(m.entries) - 1

		}
	case "pgup":
		{
			if m.cursor-listHeight > 0 {
				m.cursor -= listHeight
			} else {
				m.cursor = 0
			}
		}

	}

	return m, nil
}
