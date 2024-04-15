package database

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
	design "gosuite/design"
)

type Model struct {
	conn *db.Connection
}

func InitModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg, active bool) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View(selected bool, width int, height int) string {
	content := ""

	if m.conn == nil {
		content = "No connection"
	} else {
		content = (*m.conn).Status()
	}

	return design.CreatePane(
		1,
		"Database",
		selected,
		width,
		height,
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Render(content)),
	)

}
