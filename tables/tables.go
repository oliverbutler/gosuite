package tables

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
	design "gosuite/design"
)

type Model struct {
	Tables             []string
	SelectedTableIndex int
}

func InitModel(conn *sql.DB) Model {
	tables, err := db.GetTables(conn)
	if err != nil {
		panic(err)
	}

	return Model{
		Tables:             tables,
		SelectedTableIndex: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg, active bool, conn *sql.DB) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if active {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				if m.SelectedTableIndex > 0 && active {
					m.SelectedTableIndex--
				}
			case "down":
				if m.SelectedTableIndex < len(m.Tables)-1 && active {
					m.SelectedTableIndex++
				}
			case "enter":
				cmd = db.ExucuteSQLCmd("SELECT * FROM "+m.Tables[m.SelectedTableIndex]+" LIMIT 10", conn)
				cmds = append(cmds, cmd)

			case "i":
				cmd = db.ExucuteSQLCmd("DESCRIBE "+m.Tables[m.SelectedTableIndex], conn)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(selected bool, width int, height int) string {
	tableStyles := lipgloss.NewStyle()

	tables := make([]string, 0)

	for idx, table := range m.Tables {
		tables = append(
			tables,
			tableStyles.Foreground(design.GetBorderColor(m.SelectedTableIndex == idx && selected)).
				Render(table),
		)
	}

	content := lipgloss.NewStyle().MaxHeight(height - 15).Render(lipgloss.JoinVertical(lipgloss.Left, tables...))

	return design.CreatePane(2, "Tables", selected, width, height, content)
}
