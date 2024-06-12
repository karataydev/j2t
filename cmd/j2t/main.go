package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	//"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bTable "github.com/evertras/bubble-table/table"
	"github.com/karataymarufemre/j2t/internal/config"
	"github.com/karataymarufemre/j2t/internal/table"
)



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

	datas := []table.TableData{}	

	var f interface{}
	err = json.Unmarshal(b, &f)
	m := f.(map[string]interface{})
	for k, v := range m {
    	switch vv := v.(type) {
		case []interface{}:
			t := table.NewTableData(k)
        	for _, u := range vv {
				mm := u.(map[string]interface{})
				t.ReadJsonObj(mm)
        	}
			datas = append(datas, *t)
    	}
	}

		
	rows := []bTable.Row{}
	tables := []bTable.Model{}
	for i, v := range datas {
		rows = append(rows, bTable.NewRow(bTable.RowData{
			table.ColumnKeyId: fmt.Sprint(i),
			table.ColumnKeyTableName: v.Name(),
		}))
		tables = append(tables, *newTableModels(v))
	}

	selectTable := table.SelectTable.WithRows(rows).Focused(true)
	fmt.Println(selectTable.GetVisibleRows()[0].Data[table.ColumnKeyTableName])
	model := table.New(selectTable, tables)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	

}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func newTableModels(tab table.TableData) *bTable.Model {

	b := bTable.New(tab.Columns()).
		WithRows(tab.Rows()).
		Focused(false).
		WithMultiline(true).
		WithPageSize(30)

	return &b
}
