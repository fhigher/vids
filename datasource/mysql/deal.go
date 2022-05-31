package mysql

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("vids-mysql")

type DealMysqlData struct {
	db      *DbInfo
	tbnames []string
}

func NewDealMysqlData(db *DbInfo, tbnames string) *DealMysqlData {
	tbs := strings.Split(tbnames, ",")
	return &DealMysqlData{
		db:      db,
		tbnames: tbs,
	}
}

func (m *DealMysqlData) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, tb := range m.tbnames {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()

			result, resType, err := m.db.readTbData(t)
			if err != nil {
				log.Errorf("readTbData %s: %s", t, err)
			}

			//log.Info(t, len(result))

			m.resolveResult(result, resType)
		}(tb)
	}

	wg.Wait()
	return nil
}

func (m *DealMysqlData) resolveResult(result []interface{}, resType reflect.Type) {
	//fmt.Println(reflect.TypeOf(result[0]).Field(0).Name)
	//fmt.Println(reflect.TypeOf(result[0]).Field(0).Type)
	fmt.Println(resType)
}

/* func (m *DealMysqlData) readTbFields() (map[string][]*ColumnStruct, error) {
	if len(m.tbnames) == 0 {
		return nil, xerrors.Errorf("no specify table names")
	}

	cs, err := m.db.queryTableColumns(m.tbnames)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if len(cs) == 0 {
		return nil, xerrors.Errorf("no data found")
	}

	cmap := make(map[string][]*ColumnStruct)
	for _, c := range cs {
		if _, ok := cmap[c.TableName]; !ok {
			cmap[c.TableName] = []*ColumnStruct{c}
		} else {
			cmap[c.TableName] = append(cmap[c.TableName], c)
		}
	}

	return cmap, nil
} */
