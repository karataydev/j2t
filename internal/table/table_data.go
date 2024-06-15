package table

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/iancoleman/orderedmap"
)


type TableData struct {
	name string
	colSet orderedmap.OrderedMap
	rows []table.Row
}

func NewTableData(name string) *TableData {
	return &TableData{
		name: name,
		colSet: *orderedmap.New(),
		rows: []table.Row{},
	}
}

func (t *TableData) addCol(val string) {
	t.colSet.Set(val, true)
}

func (t *TableData) addRow(row table.Row) {
	t.rows = append(t.rows, row)
}

func (t *TableData) ReadJsonObj(m orderedmap.OrderedMap) {
	rowData := table.RowData{}
	for _, key := range m.Keys() {
		val, _ := m.Get(key)
		t.addCol(key)
		rowData[key] = readMapVal(val)
	}
	t.addRow(table.NewRow(rowData))
}

func (t TableData) ColumnsFlex() []table.Column {
	columns := []table.Column{}
	for _, k := range t.colSet.Keys() {
		columns = append(columns, table.NewFlexColumn(k, k, 1))
	}
	return columns
}

func (t TableData) Columns() []table.Column {
	columns := []table.Column{}
	for _, k := range t.colSet.Keys() {
		columns = append(columns, table.NewColumn(k, k, 30))
	}
	return columns
}


func (t TableData) Rows() []table.Row {
	return t.rows
}

func (t TableData) Name() string {
	return t.name
}

func (t *TableData) NewTableModels(isFlexCol bool) *table.Model {

	var cols []table.Column
	if isFlexCol {
		cols = t.ColumnsFlex()
	} else {
		cols = t.Columns()
	}
	b := table.New(cols).
		WithRows(t.Rows()).
		Focused(false).
		WithMultiline(true)

	return &b
}

func readMapVal(mVal interface{}) string {
	switch val := mVal.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case orderedmap.OrderedMap:
		child := NewTableData("inside")
		child.ReadJsonObj(val)
		fmt.Println(child.colSet.Keys())
		return child.NewTableModels(false).View()
	case []interface{}:
		str := []string{}
		for _, v := range val {
			str = append(str, readMapVal(v))
		}
		return lipgloss.JoinVertical(lipgloss.Left, str...)
	}
	return ""
}
