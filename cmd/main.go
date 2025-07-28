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
