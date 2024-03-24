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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.SelectedTableIndex > 0 {
				m.SelectedTableIndex--
			}
		case tea.KeyDown:
			if m.SelectedTableIndex < len(m.Tables)-1 {
				m.SelectedTableIndex++
			}
		}
	}

	return m, nil
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

	content := lipgloss.JoinVertical(lipgloss.Left, tables...)

	return design.CreatePane(2, "Tables", selected, width, height, content)
}
