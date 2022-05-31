package mysql

import "fmt"

// ColumnStruct ..
type ColumnStruct struct {
	TableName        string
	ColumnName       string
	OrdinalPosition  int
	ColumnDefault    interface{}
	IsNullable       string
	DataType         string
	MaxLength        interface{}
	NumericPrecision interface{}
	NumericScale     interface{}
	ColumnType       string
	ColumnKey        string
	Extra            string
	ColumnComment    string
}

type TableDescribe struct {
	Field     string
	FieldType string
	IsNull    string
	Key       string
	Default   interface{}
	Extra     string
}

func DebugPrint(cmap map[string][]*ColumnStruct) {
	for tb, columns := range cmap {
		fmt.Printf("TableName: %s\n", tb)
		for i, c := range columns {
			fmt.Printf("\t%d. %+v\n", i, *c)
		}
	}
}
