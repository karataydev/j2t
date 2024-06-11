package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	//"reflect"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karataymarufemre/j2t/internal/config"
)

type Table struct {
	name string
	colSet map[string]bool
	rows []table.Row
}

func NewTable(name string) *Table {
	return &Table{
		name: name,
		colSet: map[string]bool{},
		rows: []table.Row{},
	}
}

func (t *Table) addCol(val string) {
	t.colSet[val] = true
}

func (t *Table) addRow(row table.Row) {
	t.rows = append(t.rows, row)
}

func (t *Table) readJsonObj(m map[string]interface{}) {
	tRow := table.Row{}
	for key, val := range m {
		t.addCol(key)
		tRow = append(tRow, readMapVal(val))
	}
	t.addRow(tRow)
}

func (t Table) Columns() []table.Column {
	columns := []table.Column{}
	i := 0
	for k := range t.colSet {
		columns = append(columns, table.Column{Title: k, Width: max(15, len(t.rows[0][i]))})
		i++;
	}
	return columns
}

func (t Table) Rows() []table.Row {
	return t.rows
}

func readMapVal(mVal interface{}) string {
	switch val := mVal.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.2f", val)
	case []interface{}:
		str := ""
		for _, v := range val {
			str += "," + readMapVal(v)  
		}
		return str
	}
	return ""
}



func main() {
	conf, err := config.FromFlagArgs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf.FilePath)

	b, err := os.ReadFile(conf.FilePath)

    if err != nil {
        fmt.Print(err)
    }

	tables := []Table{}	

	var f interface{}
	err = json.Unmarshal(b, &f)
	m := f.(map[string]interface{})
	for k, v := range m {
    	switch vv := v.(type) {
		case []interface{}:
			t := NewTable(k)
        	for _, u := range vv {
				mm := u.(map[string]interface{})
				t.readJsonObj(mm)
        	}
			tables = append(tables, *t)
    	}
	}

	columns := []table.Column{{Title: "id", Width: 10}, {Title: "table_name", Width: 60}}
	rows := []table.Row{}
	tableModelz := []tableModels{}
	for i, v := range tables {
		rows = append(rows, table.Row{fmt.Sprint(i), v.name})
		tableModelz = append(tableModelz, *newTableModels(v))
	}


	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(42),
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

	mo := model{
		tableModels: tableModelz, 
		table: t,
		show: true,
	}
	if _, err := tea.NewProgram(mo).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	

}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))


type tableModels struct {
	table table.Model
	show bool
}

func newTableModels(tab Table) *tableModels {
	t := table.New(
		table.WithColumns(tab.Columns()),
		table.WithRows(tab.Rows()),
		table.WithFocused(true),
		table.WithHeight(42),
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

	return &tableModels{
		table: t,
	}
}
	
type model struct {
	tableModels []tableModels
	table table.Model
	show bool
	selectedIndex int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.show = true
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.show {
				m.show = false
				i, _ := strconv.Atoi(m.table.SelectedRow()[0])
				m.selectedIndex = i		
			}
		}
	}
	
	m.table, cmd = m.table.Update(msg)
	a, _ := m.tableModels[m.selectedIndex].Update(msg)
	switch tModel := a.(type) {
	case tableModels:
		m.tableModels[m.selectedIndex] = tModel
	}
	return m, cmd
}

func (m model) View() string {
	if m.show {	
		return baseStyle.Render(m.table.View()) + "\n"
	}

	return m.tableModels[m.selectedIndex].View()
}


func (m tableModels) Init() tea.Cmd { return nil }

func (m tableModels) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModels) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}