package main

import (
	"fmt"
	"log"
	"os"

	//"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bTable "github.com/evertras/bubble-table/table"
	"github.com/karataymarufemre/j2t/internal/config"
	"github.com/karataymarufemre/j2t/internal/table"
	"github.com/iancoleman/orderedmap"
)



func main() {
	conf, err := config.FromFlagArgs()
	if err != nil {
		log.Fatal(err)
	}

	b, err := os.ReadFile(conf.FilePath)

    if err != nil {
        fmt.Print(err)
    }

	datas := []table.TableData{}

	o := orderedmap.New()
	o.UnmarshalJSON(b)
	

	for _, k := range o.Keys() {
		v, _ := o.Get(k)
    	switch vv := v.(type) {
		case []interface{}:
			t := table.NewTableData(k)
        	for _, u := range vv {
				mm := u.(orderedmap.OrderedMap)
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

		tables = append(tables, v.NewTableModels(true).WithPageSize(8))
	}

	selectTable := table.SelectTable.WithRows(rows).Focused(true)
	model := table.New(selectTable, tables)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	

}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))


