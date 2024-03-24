package result

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
	design "gosuite/design"
)

type Model struct {
	table        *table.Model
	microSeconds int64
}

func InitModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg, active bool) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case db.ExecuteResult:
		columns := make([]table.Column, 0)

		for _, col := range msg.Columns {
			columns = append(columns, table.Column{
				Title: col,
				Width: 10,
			})
		}

		rows := make([]table.Row, 0)

		for _, row := range msg.Rows {
			r := make([]string, 0)

			for _, col := range msg.Columns {
				r = append(r, fmt.Sprintf("%v", row[col]))
			}

			rows = append(rows, table.Row(r))
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(10),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		t.Focus()

		m.table = &t
		m.microSeconds = msg.Microseconds
	}

	return m, nil
}

func (m Model) View(selected bool, width int, height int) string {
	content := fmt.Sprintf("Execute a query to see the results here...")

	if m.table != nil {
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			m.table.View(),
			fmt.Sprintf("Executed in %d microseconds", m.microSeconds),
		)
	}

	return design.CreatePane(4, "Results", selected, width, height, content)
}
