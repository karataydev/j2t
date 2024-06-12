package table

import (
	"github.com/evertras/bubble-table/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tables []table.Model
	selectTable table.Model
	selectedIndex int

	// Window dimensions
	totalWidth  int
	totalHeight int
}

func New(selectTable table.Model, tables []table.Model) *Model {
	return &Model{
		selectTable: selectTable,
		tables: tables,
		selectedIndex: 0,
	}
}

func (m *Model) recalculateTable() {
	m.selectTable = m.selectTable.
		WithTargetWidth(m.totalWidth)//.
		//WithMinimumHeight(m.totalHeight)
	for i,v := range m.tables {
		m.tables[i] = v.
			WithTargetWidth(m.totalWidth)//.
			//WithMinimumHeight(m.totalHeight)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)

		case "enter":
			if m.selectTable.GetFocused() {
				m.selectTable = m.selectTable.Focused(false)
				m.selectedIndex = m.selectTable.GetHighlightedRowIndex()
				m.tables[m.selectedIndex] = m.tables[m.selectedIndex].Focused(true)
			}
		case "esc":
			if !m.selectTable.GetFocused() {
				m.selectTable = m.selectTable.Focused(true)
				m.tables[m.selectedIndex] = m.tables[m.selectedIndex].Focused(false)
			}
		}
		case tea.WindowSizeMsg:
			m.totalWidth = msg.Width
			//m.totalHeight = msg.Height

			m.recalculateTable()
	}

	m.selectTable, cmd = m.selectTable.Update(msg)
	cmds = append(cmds, cmd)

	m.tables[m.selectedIndex], cmd = m.tables[m.selectedIndex].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}

func (m Model) View() string {
	if m.selectTable.GetFocused() {
		return m.selectTable.View()
	}
	return m.tables[m.selectedIndex].View()
}



