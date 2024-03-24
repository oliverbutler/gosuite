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
	tables "gosuite/tables"
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

type MainModel struct {
	db             *sql.DB
	err            error
	queryInput     string
	latestResult   dbResultModel
	terminalWidth  int
	terminalHeight int
	selectedTab    int
	tablesModel    tables.Model
}

type dbResultModel struct {
	table        table.Model
	microSeconds int64
}

type errMsg error

func initialModel() MainModel {
	conn := db.Connect()

	tablesModel := tables.InitModel(conn)

	return MainModel{
		db:          conn,
		queryInput:  "",
		err:         nil,
		selectedTab: DatabaseTab,
		tablesModel: tablesModel,
	}
}

func (m MainModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.selectedTab++

			if m.selectedTab > ResultTab {
				m.selectedTab = DatabaseTab
			}

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	if m.selectedTab == TablesTab {
		m.tablesModel, cmd = m.tablesModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) QueryView() string {
	return m.queryInput
}

func (m MainModel) ResultView() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.latestResult.table.View(),
		fmt.Sprintf("Executed in %d microseconds", m.latestResult.microSeconds),
	)
}

func getBorderColor(selected bool) lipgloss.TerminalColor {
	if selected {
		return lipgloss.Color("50")
	}
	return lipgloss.Color("255")
}

func (m MainModel) View() string {
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

	tablesTab := m.tablesModel.View(m.selectedTab == TablesTab, leftColWidth, tablesHeight)

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
