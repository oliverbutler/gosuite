package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	textarea "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// Enum for the selected tab
const (
	DatabaseTab = iota
	TablesTab
	QueryTab
	ResultTab
)

type model struct {
	db             *sql.DB
	err            error
	textarea       textarea.Model
	tables         []string
	resultTable    table.Model
	terminalWidth  int
	terminalHeight int
	selectedTab    int
}

type errMsg error

func initialModel() model {
	txt := textarea.New()
	txt.Placeholder = "Write your SQL here..."
	txt.Focus()

	conn := db.Connect()

	tables, err := db.GetTables(conn)
	if err != nil {
		panic(err)
	}

	return model{
		db:          conn,
		textarea:    txt,
		err:         nil,
		tables:      tables,
		selectedTab: DatabaseTab,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func exucuteSQL(sql string, conn *sql.DB) tea.Cmd {
	return func() tea.Msg {
		result, err := db.ExecuteSQL(conn, sql)
		if err != nil {
			return errMsg(err)
		}
		return executeSqlMsg{result: result}
	}
}

type executeSqlMsg struct {
	result db.ExecuteResult
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
	case executeSqlMsg:

		columns := make([]table.Column, 0)

		for _, col := range msg.result.Columns {
			columns = append(columns, table.Column{
				Title: col,
				Width: 10,
			})
		}

		rows := make([]table.Row, 0)

		for _, row := range msg.result.Rows {
			r := make([]string, 0)

			for _, col := range msg.result.Columns {
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

		m.resultTable = t

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			if !m.textarea.Focused() {

				cmd = exucuteSQL(m.textarea.Value(), m.db)
				cmds = append(cmds, cmd)
			} else {
				m.textarea, cmd = m.textarea.Update(msg)
			}

		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}

		case tea.KeyCtrlC:
			return m, tea.Quit

		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
			} else {
				m.textarea, cmd = m.textarea.Update(msg)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// stripANSI removes ANSI escape sequences from a string.
func stripANSI(str string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(str, "")
}

// 'selected' determines whether the text is bolded.
func addTextToBorder(content, index, text string, selected bool) string {
	lines := strings.Split(content, "\n")
	if len(lines) < 1 {
		return content // If there's no content, just return it unchanged.
	}

	// Render the insertion text with optional bold styling.
	insertionText := lipgloss.NewStyle().
		Bold(selected).
		Foreground(getBorderColor(selected)).
		Render("[" + index + "] " + text)
	// Calculate the visible length of the insertionText by stripping ANSI codes.
	visibleInsertionLength := len(stripANSI(insertionText))

	// Calculate insertion points based on the actual visible content length.
	insertAt := 3 // Assumes a fixed starting position for the insertion.

	// Convert the first line to runes to handle potential Unicode characters correctly.
	runes := []rune(lines[0])

	// Calculate the cut index based on visible content, adjusting for the border's visual length.
	cutIndex := insertAt + visibleInsertionLength
	if cutIndex > len(runes) {
		cutIndex = len(runes) // Ensure cutIndex does not exceed the line length.
	}

	// Reassemble the modified top border.
	beforeInsertion := string(runes[:insertAt])
	afterInsertion := string(runes[cutIndex:])

	// Construct the new first line with the inserted text.
	lines[0] = beforeInsertion + insertionText + afterInsertion

	return strings.Join(lines, "\n")
}

func (m model) TablesView() string {
	tableStyles := lipgloss.NewStyle()

	tables := make([]string, 0)

	for _, table := range m.tables {
		tables = append(tables, tableStyles.Render(table))
	}

	return lipgloss.JoinVertical(lipgloss.Left, tables...)
}

func (m model) QueryView() string {
	return m.textarea.View()
}

func (m model) ResultView() string {
	return m.resultTable.View()
}

func getBorderColor(selected bool) lipgloss.TerminalColor {
	if selected {
		return lipgloss.Color("50")
	}
	return lipgloss.Color("255")
}

func (m model) View() string {
	// Tab is a bordered box, with a name near the top left with a number e.g. 1. Database or 2. Tables
	tabStyles := lipgloss.NewStyle().
		Padding(1, 1).
		Border(lipgloss.RoundedBorder()).
		MarginRight(1)

	safeWidth := m.terminalWidth - 5
	safeHeight := m.terminalHeight - 5

	leftColWidth := 40
	rightColWidth := safeWidth - leftColWidth

	databaseHeight := 5
	tablesHeight := safeHeight - databaseHeight

	queryHeight := 10
	resultHeight := safeHeight - queryHeight

	databaseTab := addTextToBorder(
		tabStyles.BorderForeground(getBorderColor(m.selectedTab == DatabaseTab)).
			Width(leftColWidth).
			Height(databaseHeight).
			Render("127.0.0.1:3306 (LOCAL)"),
		"1",
		"Database",
		m.selectedTab == DatabaseTab,
	)
	tablesTab := addTextToBorder(
		tabStyles.
			BorderForeground(
				getBorderColor(m.selectedTab == TablesTab),
			).Width(leftColWidth).Height(tablesHeight).Render(m.TablesView()),
		"2",
		"Tables",
		m.selectedTab == TablesTab,
	)
	queryTab := addTextToBorder(
		tabStyles.
			BorderForeground(
				getBorderColor(m.selectedTab == QueryTab),
			).Width(rightColWidth).Height(queryHeight).Render(m.QueryView()),
		"3",
		"Query",
		m.selectedTab == QueryTab,
	)
	resultTab := addTextToBorder(
		tabStyles.BorderForeground(getBorderColor(m.selectedTab == ResultTab)).
			Width(rightColWidth).
			Height(resultHeight).
			Render(m.ResultView()),
		"4",
		"Result",
		m.selectedTab == ResultTab,
	)

	leftCol := lipgloss.JoinVertical(lipgloss.Left, databaseTab, tablesTab)
	rightCol := lipgloss.JoinVertical(lipgloss.Left, queryTab, resultTab)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	return layout
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
