package main

import (
	"database/sql"
	"fmt"
	"os"

	textarea "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
	design "gosuite/design"
	result "gosuite/result"
	tables "gosuite/tables"
)

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
	terminalWidth  int
	terminalHeight int
	selectedTab    int
	tablesModel    tables.Model
	resultModel    result.Model
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

	m.tablesModel, cmd = m.tablesModel.Update(msg, m.selectedTab == TablesTab, m.db)
	cmds = append(cmds, cmd)

	m.resultModel, cmd = m.resultModel.Update(msg, m.selectedTab == ResultTab)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MainModel) QueryView() string {
	return m.queryInput
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

	resultTab := m.resultModel.View(m.selectedTab == ResultTab, rightColWidth, resultHeight)

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
