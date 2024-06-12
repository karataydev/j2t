package table

import (
	"github.com/evertras/bubble-table/table"
)

const (
	ColumnKeyId       	= "id"
	ColumnKeyTableName  = "table name"
)

var SelectTable = table.New([]table.Column{
	table.NewColumn(ColumnKeyId, "id", 10),
	table.NewFlexColumn(ColumnKeyTableName, "table name", 1),
})
