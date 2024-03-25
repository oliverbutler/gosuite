package query

import (
	"database/sql"
	"regexp"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
	design "gosuite/design"
)

var (
	keywordStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("140"))
	stringStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	commentStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	functionStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("198"))
	operatorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("198"))
	identifierStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("87"))
)

func sqlHighlighter(sql string) string {
	// Define regular expressions for different SQL syntax elements
	keywordRegex := regexp.MustCompile(
		`\b(?i)(SELECT|FROM|WHERE|JOIN|GROUP BY|ORDER BY|LIMIT|OFFSET|INSERT INTO|UPDATE|DELETE|CREATE|DROP|ALTER|TABLE|VIEW|INDEX|TRIGGER|PROCEDURE|FUNCTION)\b`,
	)
	stringRegex := regexp.MustCompile(`'[^']*'`)
	commentRegex := regexp.MustCompile(`--.*`)
	functionRegex := regexp.MustCompile(
		`\b(?i)(COUNT|SUM|AVG|MIN|MAX|CONCAT|SUBSTRING|TRIM|LENGTH|UPPER|LOWER|ROUND|COALESCE)\b`,
	)
	operatorRegex := regexp.MustCompile(
		`(?i)(\b(AND|OR|NOT|LIKE|BETWEEN|IN|IS NULL|EXISTS)\b|[=<>!]+)`,
	)
	identifierRegex := regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*\b`)

	// Apply syntax highlighting using Lipgloss styles
	sql = keywordRegex.ReplaceAllStringFunc(sql, func(match string) string {
		return keywordStyle.Render(match)
	})
	sql = stringRegex.ReplaceAllStringFunc(sql, func(match string) string {
		return stringStyle.Render(match)
	})
	sql = commentRegex.ReplaceAllStringFunc(sql, func(match string) string {
		return commentStyle.Render(match)
	})
	sql = functionRegex.ReplaceAllStringFunc(sql, func(match string) string {
		return functionStyle.Render(match)
	})
	sql = operatorRegex.ReplaceAllStringFunc(sql, func(match string) string {
		return operatorStyle.Render(match)
	})
	sql = identifierRegex.ReplaceAllStringFunc(sql, func(match string) string {
		return identifierStyle.Render(match)
	})

	return sql
}

type Model struct {
	Input textarea.Model
}

func InitModel() Model {
	ta := textarea.New()

	// Remove the white line on the left
	ta.ShowLineNumbers = true
	ta.Prompt = ""

	return Model{
		Input: ta,
	}
}

func FocusOnQuery() tea.Cmd {
	return func() tea.Msg {
		return FocusOnQueryMsg{}
	}
}

type FocusOnQueryMsg struct{}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg, active bool, conn *sql.DB) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case db.ExecuteResult:
		m.Input.SetValue(msg.Query)
	case FocusOnQueryMsg:
		m.Input.Focus()
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if active && !m.Input.Focused() {
				cmd = db.ExucuteSQLCmd(m.Input.Value(), conn)
				cmds = append(cmds, cmd)
			}
		case "esc":
			if m.Input.Focused() {
				m.Input.Blur()
			}

		default:
			if active {
				m.Input, cmd = m.Input.Update(msg)
			}

		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View(selected bool, width int, height int) string {
	return design.CreatePane(3, "Query", selected, width, height, sqlHighlighter(m.Input.View()))
}
