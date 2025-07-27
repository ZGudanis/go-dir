package godir

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	colorGreen      = lipgloss.Color("#1fb009")
	colorWhite      = lipgloss.Color("#dededeff")
	colorLightGreen = lipgloss.Color("#8df87dff")
	colorBlue       = lipgloss.Color("#00166cff")
)

var (
	dirListWidth float32 = 0.3
	width        int
	height       int

	listStyle     lipgloss.Style
	previewStyle  lipgloss.Style
	mainViewStyle lipgloss.Style
	folderStyle   lipgloss.Style = lipgloss.NewStyle().Foreground(colorGreen)
	selectedStyle lipgloss.Style = lipgloss.NewStyle().Foreground(colorWhite).Background(colorBlue)

	subtle = lipgloss.AdaptiveColor{Light: "#d9dccf", Dark: "#383838"}

	text textarea.Model
)

func init() {
	width, height, _ = term.GetSize(int(os.Stdout.Fd()))

	listStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false, false).
		BorderForeground(subtle).
		MarginRight(0).
		MarginBottom(0).
		Height(verticalFill(.7)).
		Width(horizontalFill(dirListWidth)).
		MaxHeight(verticalFill(.9))

	previewStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false, false).
		BorderForeground(subtle).
		MarginRight(0).
		MarginBottom(0).
		Height(verticalFill(.9)).
		MaxHeight(verticalFill(.9))

	mainViewStyle = lipgloss.NewStyle().
		MaxHeight(height).
		MaxWidth(width)

	text = textarea.New()
	text.SetWidth(horizontalFill(1-dirListWidth) - 7)
	text.SetHeight(verticalFill(.9))
}

func RenderList(model DirectoryModel) string {
	text.SetValue(preview(model))

	main := lipgloss.JoinHorizontal(lipgloss.Top,
		listStyle.Render(lipgloss.JoinVertical(lipgloss.Left, parseEntryNames(model, listStyle.GetHeight())...)),
		// previewStyle.Render(strconv.Itoa(model.cursor)))
		text.View())

	return mainViewStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
		fmt.Sprintf("%s\n", model.cwd),
		main,
		"\nPress q to exit\n"))
}

func parseEntryNames(model DirectoryModel, count int) []string {
	names := make([]string, len(model.entries))
	var iterationEnd int = 0
	var iterator int = 0
	if len(model.entries) < count {
		iterationEnd = len(model.entries)
	} else {
		iterationEnd = min(model.cursor+count, len(model.entries))
		iterator = iterationEnd - count
	}

	for i, entry := range model.entries[iterator:iterationEnd] {
		if entry.IsDir() {
			names[i] = folder(entry.Name(), i+iterator == model.cursor)
			continue
		}
		if i+iterator == model.cursor {
			names[i] = selectedStyle.Render(entry.Name())
			continue
		}
		names[i] = entry.Name()
	}
	return names
}

var folder = func(s string, selected bool) string {
	var style = folderStyle
	if selected {
		style = style.Inherit(selectedStyle)
	}

	return style.Render(s)
}

func preview(model DirectoryModel) string {
	if len(model.entries) == 0 {
		return ""
	}

	entry := model.entries[model.cursor]

	if !entry.Type().IsDir() {
		buffer := make([]byte, 2048)
		file, err := os.Open(model.cwd + "/" + entry.Name())
		if err != nil {
			return err.Error()
		}
		defer file.Close()

		read, err := file.Read(buffer)
		if err != nil {
			return "Unable to read"
		}
		return string(buffer[0:read])
	} else {
		names := strings.Builder{}
		dir, err := os.ReadDir(model.cwd + "/" + entry.Name())
		if err != nil {
			return err.Error()
		}
		for _, entry := range dir {
			// if entry.IsDir() {
			// 	names.WriteString(folder(strings.Trim(entry.Name(), "")+"\n", false))
			// 	continue
			// }
			names.WriteString(strings.Trim(entry.Name(), "") + "\n")
		}

		return names.String()
	}
}

func verticalFill(fill float32) int {
	return int(float32(height) * fill)
}
func horizontalFill(fill float32) int {
	return int(float32(width) * fill)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
