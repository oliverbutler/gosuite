package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	textarea "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
	design "gosuite/design"
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
	db                 *sql.DB
	err                error
	textarea           textarea.Model
	tables             []string
	resultTable        table.Model
	terminalWidth      int
	terminalHeight     int
	selectedTab        int
	selectedTableIndex int
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
		db:                 conn,
		textarea:           txt,
		err:                nil,
		tables:             tables,
		selectedTab:        DatabaseTab,
		selectedTableIndex: 0,
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

	if m.selectedTab == TablesTab {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyDown:
				m.selectedTableIndex++
				if m.selectedTableIndex > len(m.tables) {
					m.selectedTableIndex = 0
				}

			case tea.KeyUp:
				m.selectedTableIndex--
				if m.selectedTableIndex < 0 {
					m.selectedTableIndex = len(m.tables) - 1
				}
			}
		}
	}

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

		case tea.KeyTab:
			m.selectedTab++
			if m.selectedTab > ResultTab {
				m.selectedTab = DatabaseTab
			}

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

func (m model) TablesView() string {
	tableStyles := lipgloss.NewStyle()

	tables := make([]string, 0)

	for idx, table := range m.tables {
		tables = append(
			tables,
			tableStyles.Foreground(getBorderColor(m.selectedTableIndex == idx)).Render(table),
		)
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
	safeWidth := m.terminalWidth - 5
	safeHeight := m.terminalHeight - 5

	leftColWidth := 40
	rightColWidth := safeWidth - leftColWidth

	databaseHeight := 5
	tablesHeight := safeHeight - databaseHeight

	queryHeight := 10
	resultHeight := safeHeight - queryHeight

	databaseTab := design.CreatePane(
		1,
		"Database",
		m.selectedTab == DatabaseTab,
		leftColWidth,
		databaseHeight,
		"127.0.0.1:3306",
	)

	tablesTab := design.CreatePane(
		2,
		"Tables",
		m.selectedTab == TablesTab,
		leftColWidth,
		tablesHeight,
		m.TablesView(),
	)

	queryTab := design.CreatePane(
		3,
		"Query",
		m.selectedTab == QueryTab,
		rightColWidth,
		queryHeight,
		m.QueryView(),
	)

	resultTab := design.CreatePane(
		4,
		"Result",
		m.selectedTab == ResultTab,
		rightColWidth,
		resultHeight,
		m.ResultView(),
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
