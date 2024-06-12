package table

import (
	"strconv"
	"github.com/evertras/bubble-table/table"
)


type TableData struct {
	name string
	colSet map[string]bool
	rows []table.Row
}

func NewTableData(name string) *TableData {
	return &TableData{
		name: name,
		colSet: map[string]bool{},
		rows: []table.Row{},
	}
}

func (t *TableData) addCol(val string) {
	t.colSet[val] = true
}

func (t *TableData) addRow(row table.Row) {
	t.rows = append(t.rows, row)
}

func (t *TableData) ReadJsonObj(m map[string]interface{}) {
	rowData := table.RowData{}
	for key, val := range m {
		t.addCol(key)
		rowData[key] = readMapVal(val)
	}
	t.addRow(table.NewRow(rowData))
}

func (t TableData) Columns() []table.Column {
	columns := []table.Column{}
	i := 0
	for k := range t.colSet {
		columns = append(columns, table.NewFlexColumn(k, k, 1))
		i++;
	}
	return columns
}

func (t TableData) Rows() []table.Row {
	return t.rows
}

func (t TableData) Name() string {
	return t.name
}

func readMapVal(mVal interface{}) string {
	switch val := mVal.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case []interface{}:
		str := ""
		for i, v := range val {
			str += readMapVal(v)
			if i < len(val) - 1 {
				str += "\n"
			}
		}
		return str
	}
	return ""
}
