package result

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
	design "gosuite/design"
)

type Cursor struct {
	Row    int
	Column int
}

type Model struct {
	result *db.ExecuteResult
	cursor Cursor
}

func InitModel() Model {
	return Model{
		cursor: Cursor{0, 0},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg, active bool) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case db.ExecuteResult:
		m.result = &msg
	case tea.KeyMsg:
		if !active {
			return m, nil
		}

		switch msg.String() {
		case "up":
			if m.cursor.Row > 0 {
				m.cursor.Row--
			}
		case "down":
			if m.cursor.Row < len(m.result.Rows)-1 {
				m.cursor.Row++
			}
		case "left":
			if m.cursor.Column > 0 {
				m.cursor.Column--
			}
		case "right":
			if m.cursor.Column < len(m.result.Columns)-1 {
				m.cursor.Column++
			}

		}
	}

	return m, nil
}

func getWidthFromColumn(column string) int {
	return len(column) + 10
}

func truncateString(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	return s
}

func renderColumns(result *db.ExecuteResult) string {
	content := make([]string, 0)

	for _, column := range result.Columns {
		content = append(content,
			lipgloss.NewStyle().
				Padding(0, 1).
				Background(lipgloss.Color("238")).
				Width(getWidthFromColumn(column)).
				Render(fmt.Sprintf("%v", column)))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, content...)
}

func renderRow(columns []string, data map[string]interface{}, cursorColumnIndex int) string {
	var content string

	for idx, column := range columns {

		width := getWidthFromColumn(column)
		truncated := truncateString(fmt.Sprintf("%v", data[column]), width-2)

		if cursorColumnIndex == idx {
			content += lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#ffffff")).
				Padding(0, 1).
				Width(width).
				Render(truncated)
		} else {
			content += lipgloss.NewStyle().
				Padding(0, 1).
				Width(width).
				Render(truncated)
		}
	}

	return content
}

func renderRows(result *db.ExecuteResult, cursor Cursor) string {
	var content string

	for idx, row := range result.Rows {

		isRowSelected := cursor.Row == idx

		columnIndex := -1

		if isRowSelected {
			columnIndex = cursor.Column
		}

		content += renderRow(result.Columns, row, columnIndex) + "\n"
	}

	return content
}

func (m Model) View(selected bool, width int, height int) string {
	content := fmt.Sprintf("Execute a query to see the results here...")

	if m.result != nil {
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			renderColumns(m.result),
			renderRows(m.result, m.cursor),
			fmt.Sprintf("Executed in %d microseconds", m.result.Microseconds),
		)
	}

	return design.CreatePane(4, "Results", selected, width, height, content)
}
