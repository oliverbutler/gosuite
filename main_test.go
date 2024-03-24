package main

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestView(t *testing.T) {
	m := initialModel()

	res := m.View()

	println(res)
}

func TestBorderStyling(t *testing.T) {
	res := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("24")).
		Render("Hello, World!")

	println(res)
}
