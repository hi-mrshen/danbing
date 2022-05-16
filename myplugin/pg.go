package myplugin

import (
	"danbing/plugin"
	"danbing/task"
	"encoding/json"
	"fmt"

	"database/sql"

	_ "github.com/lib/pq"
)

type PgReader struct {
	Query   *task.Query
	Connect *task.Connect
	db      *sql.DB
}

func (reader *PgReader) Init(tq *task.Query, tc *task.Connect) {
	var pool *sql.DB
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		tc.Host, tc.Port, tc.Username, tc.Password, tc.Database)
	pool, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Println(err)
	}
	err = pool.Ping()
	if err != nil {
		panic(err)
	}
	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)
	reader.db = pool
	reader.Query = tq
}
func (reader *PgReader) Name() string {
	return "pgsqlreader"
}

func (reader *PgReader) Split(taskNum int) []plugin.ReaderPlugin {
	plugins := make([]plugin.ReaderPlugin, 0)
	for i := 0; i < taskNum; i++ {
		reader.Query.Offset = i * reader.Query.Size
		plugins = append(plugins, reader)
	}
	return plugins
}

func (reader *PgReader) Reader() string {
	result := make([]map[string]string, 0)

	rows, err := reader.db.Query(reader.Query.SQL)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	if len(cols) > 0 {
		for rows.Next() {
			buff := make([]interface{}, len(cols))
			data := make([][]byte, len(cols)) //数据库中的NULL值可以扫描到字节中
			for i, _ := range buff {
				buff[i] = &data[i]
			}
			rows.Scan(buff...)
			dataKv := make(map[string]string)
			for k, col := range data {
				dataKv[cols[k]] = string(col)
			}
			result = append(result, dataKv)
		}
	}
	s, _ := json.Marshal(result)
	return string(s)
}

func (reader *PgReader) Close() {
	reader.db.Close()
}

// TODO: init必须手动维护
func init() {
	plugin.Register(&PgReader{})
}