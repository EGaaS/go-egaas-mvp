package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/EGaaS/go-egaas-mvp/packages/model"
	_ "github.com/lib/pq"
)

type dbWorker struct {
	*sql.DB
	config map[string]string
}

var (
	CONNECTION_STRING_ERROR = errors.New("Connection string error")
	CONNECTION_ERROR        = errors.New("Can't connect to database")
)

func Connect(config map[string]string) (*dbWorker, error) {
	if len(config["db_user"]) == 0 || len(config["db_password"]) == 0 || len(config["db_name"]) == 0 {
		return &dbWorker{}, CONNECTION_STRING_ERROR
	}

	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=%s",
			config["db_user"], config["db_password"], config["db_name"], config["db_port"]))
	if err != nil || db.Ping() != nil {
		return &dbWorker{}, CONNECTION_ERROR
	}
	return &dbWorker{db, config}, nil
}

func (db *dbWorker) Run(m *model.Model) {
	// in this case the query is SELECT
	if m.ReturnValue != nil && len(m.ReturnValue) > 0 {
		rows, err := db.Query(m.Query())
		if err != nil {
			m.Error = err
			return
		}
		defer rows.Close()

		values := make([][]byte, len(m.ReturnValue))
		scanArgs := make([]interface{}, len(m.ReturnValue))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				m.Error = err
				return
			}

			for i, col := range values {
				if col == nil {
					m.ReturnValue[i].SetNull(true)
				} else {
					m.ReturnValue[i].FromBytes(col)
				}
			}
		}
		// and in this case the query is not SELECT
	} else {
		_, err := db.Exec(m.Query())
		if err != nil {
			m.Error = err
			return
		}
	}
}
