package godir

import (
	"os"
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
			{
				return updatePath(m, Outof), nil
			}
		}
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
