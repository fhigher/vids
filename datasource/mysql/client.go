package mysql

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DbInfo struct {
	host, port, user, pass, dbname string
	maxConns                       int
	dbObj                          *sql.DB
}

func NewDbInfo(host, port, user, pass, dbname string, maxConns int) *DbInfo {
	return &DbInfo{
		host,
		port,
		user,
		pass,
		dbname,
		maxConns,
		nil,
	}
}

func (db *DbInfo) InitDB() {
	var err error

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s&parseTime=true", db.user, db.pass, db.host,
		db.port, db.dbname, url.QueryEscape("Asia/Shanghai"))
	dbObj, err := sql.Open("mysql", dns)

	if nil != err {
		panic(err)
	}
	dbObj.SetMaxOpenConns(db.maxConns)
	// 空闲链接数可以为最大链接数的百分比配置
	dbObj.SetMaxIdleConns(10)
	dbObj.SetConnMaxIdleTime(5 * time.Minute)


	err = dbObj.Ping()
	if nil != err {
		panic(err)
	}

	db.dbObj = dbObj
}

/* func (db *DbInfo) queryTableColumns(tbnames []string) ([]*ColumnStruct, error) {

	errstr := "QueryTableColumns"

	fields := "table_name, column_name, ordinal_position, column_default, is_nullable, data_type," +
		" character_maximum_length, numeric_precision, numeric_scale, column_type, column_key, extra, column_comment"

	tbs := ""
	for _, s := range tbnames {
		tbs += fmt.Sprintf("'%s',", s)
	}

	conditions := fmt.Sprintf("table_name in (%s)", strings.TrimRight(tbs, ","))

	sqlstr := fmt.Sprintf("SELECT %s FROM information_schema.`COLUMNS` WHERE %s ", fields, conditions)
	log.Debug(sqlstr)

	rows, err := db.dbObj.Query(sqlstr)
	if nil != err {
		return nil, fmt.Errorf("%s.Query: %s", errstr, err.Error())
	}

	defer rows.Close()

	cs := make([]*ColumnStruct, 0)

	for rows.Next() {
		c := ColumnStruct{}
		err := rows.Scan(&c.TableName, &c.ColumnName, &c.OrdinalPosition, &c.ColumnDefault, &c.IsNullable, &c.DataType, &c.MaxLength,
			&c.NumericPrecision, &c.NumericScale, &c.ColumnType, &c.ColumnKey, &c.Extra, &c.ColumnComment)

		if err != nil {
			return cs, fmt.Errorf("%s.Scan: %s", errstr, err.Error())
		}
		cs = append(cs, &c)
	}

	return cs, nil
} */

func (db *DbInfo) readTbData(tb string) ([]interface{}, reflect.Type, error) {
	errstr := "readTbData"

	sqlstr := fmt.Sprintf("select * from %s limit 2", tb)

	rows, err := db.dbObj.Query(sqlstr)
	if err != nil {
		return nil, nil, fmt.Errorf("%s.Query: %s", errstr, err.Error())
	}

	defer rows.Close()

	dd, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, fmt.Errorf("%s.ColumnTypes: %s", errstr, err.Error())
	}
	//fmt.Printf("%s, %+v\n", dd[0].Name(), dd[0].ScanType())
	t := []reflect.StructField{}
	for i := 0; i < len(dd); i++ {
		name := dd[i].Name()
		s := strings.ToUpper(name[:1]) + name[1:]
		tt := reflect.StructField{
			Name: s,
			Type: dd[i].ScanType(),
		}

		t = append(t, tt)
	}

	st := reflect.StructOf(t)
	result := []interface{}{}
	for rows.Next() {
		v := reflect.New(st).Elem()
		args := []interface{}{}
		for i := 0; i < v.NumField(); i++ {
			args = append(args, v.Field(i).Addr().Interface())
		}
		err = rows.Scan(args...)
		if err != nil {
			return nil, nil, fmt.Errorf("%s.Scan: %s", errstr, err.Error())
		}

		result = append(result, v.Interface())
		//fmt.Printf("%+v\n", v.Interface())
	}

	return result, st, nil
}
