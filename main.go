package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "root"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	c, _ := db.Conn(context.Background())

	type Req struct {
		SQL string `json:"sql"`
	}
	type Column struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	type Res struct {
		Columns []Column `json:"columns"`
		Data    []any    `json:"data"`
	}
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Req
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			panic(err)
		}

		stm, err := c.QueryContext(r.Context(), req.SQL)
		if err != nil {
			panic(err)
		}

		cols, err := stm.ColumnTypes()
		if err != nil {
			panic(err)
		}
		nCols := len(cols)

		var rows []any
		for stm.Next() {
			rowColsPtrs := make([]*any, nCols)
			for i, colType := range cols {
				val := reflect.New(colType.ScanType())
				ptr := val.Interface()
				rowColsPtrs[i] = &ptr
			}

			// convert []*any into []any to make go happy
			dest := make([]any, len(rowColsPtrs))
			for i, v := range rowColsPtrs {
				dest[i] = v
			}

			err = stm.Scan(dest...)
			if err != nil {
				panic(err)
			}
			rows = append(rows, rowColsPtrs)
		}

		colsRes := make([]Column, 0, nCols)
		for _, col := range cols {
			colsRes = append(colsRes, Column{
				Name: col.Name(),
				Type: col.DatabaseTypeName(),
			})
		}

		json.NewEncoder(w).Encode(Res{
			Columns: colsRes,
			Data:    rows,
		})
	}))

	fmt.Println("Successfully connected!")
}
