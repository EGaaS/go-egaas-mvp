package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/EGaaS/go-egaas-mvp/packages/dal/model"
	_ "github.com/lib/pq"
)

type pgWorker struct {
	*sql.DB
	config map[string]string
}

var (
	CONNECTION_STRING_ERROR = errors.New("Connection string error")
	CONNECTION_ERROR        = errors.New("Can't connect to database")
)

func PgConnect(login, pass, dbName string, dbPort int, config map[string]string) (*pgWorker, error) {
	if len(login) == 0 || len(pass) == 0 || len(dbName) == 0 {
		return &pgWorker{}, CONNECTION_STRING_ERROR
	}

	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=%s",
			login, pass, dbName, dbPort))
	if err != nil || db.Ping() != nil {
		return &pgWorker{}, CONNECTION_ERROR
	}
	return &pgWorker{db, config}, nil
}

func (db *pgWorker) Run(m *model.Model) {
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
