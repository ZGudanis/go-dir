package godir

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	width       = 96
	columnWidth = 30

	colorGreen      = lipgloss.Color("#1fb009")
	colorWhite      = lipgloss.Color("#dededeff")
	colorLightGreen = lipgloss.Color("#8df87dff")
	colorBlue       = lipgloss.Color("#00166cff")
)

func RenderList(model DirectoryModel) string {
	list := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false, false).
		BorderForeground(subtle).
		MarginRight(2).
		Height(8).
		Width(columnWidth + 1)

	return list.Render(
		lipgloss.JoinVertical(lipgloss.Left, parseEntryNames(model)...),
	)
}

var subtle = lipgloss.AdaptiveColor{Light: "#d9dccf", Dark: "#383838"}

func parseEntryNames(model DirectoryModel) []string {
	names := make([]string, len(model.entries))
	for i, entry := range model.entries {
		if entry.IsDir() {
			names[i] = folder(entry.Name(), i == model.cursor)
			continue
		}
		if i == model.cursor {
			names[i] = selectedStyle.Render(entry.Name())
			continue
		}
		names[i] = entry.Name()
	}
	return names
}

var folderStyle = lipgloss.NewStyle().Foreground(colorGreen)
var selectedStyle = lipgloss.NewStyle().Foreground(colorWhite).Background(colorBlue)

var folder = func(s string, selected bool) string {
	var style = folderStyle
	if selected {
		style = style.Inherit(selectedStyle)
	}

	return style.Render(s)
}
