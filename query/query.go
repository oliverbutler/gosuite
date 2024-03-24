package query

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"

	db "gosuite/db"
	design "gosuite/design"
)

type Model struct {
	query string
}

func InitModel() Model {
	return Model{
		query: "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg, active bool, conn *sql.DB) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case db.ExecuteResult:
		m.query = msg.Query
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if active {
				cmd = db.ExucuteSQLCmd(m.query, conn)
				cmds = append(cmds, cmd)
			}
		case tea.KeyBackspace:
			if active && len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
			}

		case tea.KeyTab:

		default:
			if active {
				m.query += msg.String()
			}

		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(selected bool, width int, height int) string {
	return design.CreatePane(3, "Query", selected, width, height, m.query)
}
