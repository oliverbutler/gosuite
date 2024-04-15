package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	textarea "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	db "gosuite/db"
	query "gosuite/query"
	result "gosuite/result"
	"gosuite/services/config"
	tables "gosuite/tables"
	database "gosuite/views/database"
)

const (
	DatabaseTab = iota
	TablesTab
	QueryTab
	ResultTab
)

type MainModel struct {
	config *config.AppConfig

	connection db.Connection

	err            error
	terminalWidth  int
	terminalHeight int
	selectedTab    int

	databaseModel database.Model
	tablesModel   tables.Model
	resultModel   result.Model
	queryModel    query.Model

	// Keys
	keys keyMap
	help help.Model
}

type keyMap map[string]key.Binding

var QuitKey = key.NewBinding(
	key.WithKeys("q", "ctrl+c"),
	key.WithHelp("q", "Quit"),
)

var TabKey = key.NewBinding(
	key.WithKeys("tab"),
	key.WithHelp("tab", "Next tab"),
)

var ShiftTabKey = key.NewBinding(
	key.WithKeys("shift+tab"),
	key.WithHelp("shift+tab", "Previous tab"),
)

var FocusQueryKey = key.NewBinding(
	key.WithKeys("/"),
	key.WithHelp("/", "Focus on query"),
)

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		QuitKey,
		TabKey,
		ShiftTabKey,
		FocusQueryKey,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			QuitKey,
			TabKey,
			ShiftTabKey,
			FocusQueryKey,
		},
	}
}

type errMsg error

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

	case query.FocusOnQueryMsg:
		m.selectedTab = QueryTab

	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab":
			m.selectedTab--

			if m.selectedTab < DatabaseTab {
				m.selectedTab = ResultTab
			}

		case "tab":

			m.selectedTab++

			if m.selectedTab > ResultTab {
				m.selectedTab = DatabaseTab
			}

		case "1":
			if !m.queryModel.Input.Focused() {
				m.selectedTab = DatabaseTab
			}
		case "2":
			if !m.queryModel.Input.Focused() {
				m.selectedTab = TablesTab
			}
		case "3":
			if !m.queryModel.Input.Focused() {
				m.selectedTab = QueryTab
			}
		case "4":
			if !m.queryModel.Input.Focused() {
				m.selectedTab = ResultTab
			}

		case "/":
			cmd = query.FocusOnQuery()
			cmds = append(cmds, cmd)

		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		}
	}

	m.databaseModel, cmd = m.databaseModel.Update(msg, m.selectedTab == DatabaseTab)
	cmds = append(cmds, cmd)

	m.tablesModel, cmd = m.tablesModel.Update(msg, m.selectedTab == TablesTab, &m.connection)
	cmds = append(cmds, cmd)

	m.resultModel, cmd = m.resultModel.Update(msg, m.selectedTab == ResultTab)
	cmds = append(cmds, cmd)

	m.queryModel, cmd = m.queryModel.Update(msg, m.selectedTab == QueryTab, &m.connection)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	safeWidth := m.terminalWidth - 5
	safeHeight := m.terminalHeight - 6

	leftColWidth := 25
	rightColWidth := safeWidth - leftColWidth

	databaseHeight := 5
	tablesHeight := safeHeight - databaseHeight

	queryHeight := 10
	resultHeight := safeHeight - queryHeight

	databaseTab := m.databaseModel.View(m.selectedTab == DatabaseTab, leftColWidth, databaseHeight)
	tablesTab := m.tablesModel.View(m.selectedTab == TablesTab, leftColWidth, tablesHeight)
	queryTab := m.queryModel.View(m.selectedTab == QueryTab, rightColWidth, queryHeight)
	resultTab := m.resultModel.View(m.selectedTab == ResultTab, rightColWidth, resultHeight)

	leftCol := lipgloss.JoinVertical(lipgloss.Left, databaseTab, tablesTab)
	rightCol := lipgloss.JoinVertical(lipgloss.Left, queryTab, resultTab)

	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol),
		m.help.View(m.keys),
	)

	return layout
}

func main() {
	config, err := config.GetConfig()

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	databasesLength := len(config.Databases)

	if databasesLength == 0 {
		fmt.Printf("No databases found in config file, please add at least one database.")
		os.Exit(1)
	}

	firstConfig := config.Databases[0]

	conn := db.ConnectionPending{
		Config: &firstConfig,
	}

	databaseModel := database.InitModel()
	tablesModel := tables.InitModel()
	resultModel := result.InitModel()
	queryModel := query.InitModel()

	model := MainModel{
		connection:    conn,
		err:           nil,
		selectedTab:   TablesTab,
		databaseModel: databaseModel,
		tablesModel:   tablesModel,
		resultModel:   resultModel,
		queryModel:    queryModel,
		keys: keyMap{
			"Quit":       QuitKey,
			"Tab":        TabKey,
			"ShiftTab":   ShiftTabKey,
			"FocusQuery": FocusQueryKey,
		},
		help: help.NewModel(),
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
