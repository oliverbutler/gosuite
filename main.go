package main

import (
	"database/sql"
	"fmt"
	"os"

	textarea "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	db "gosuite/db"
)

type model struct {
	db       *sql.DB
	err      error
	textarea textarea.Model
	result   string
}

type errMsg error

func initialModel() model {
	txt := textarea.New()
	txt.Placeholder = "Write your SQL here..."
	txt.Focus()

	db := db.Connect()

	return model{
		db:       db,
		textarea: txt,
		err:      nil,
		result:   "",
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

		fmt.Println(result)

		return executeSqlMsg{sql: "soe string"}
	}
}

type executeSqlMsg struct {
	sql string
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case executeSqlMsg:
		m.result = msg.sql
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

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		m.textarea.View(),
		m.result,
		"(ctrl+c to quit)",
	) + "\n\n"
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
