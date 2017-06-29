package model

import "github.com/EGaaS/go-egaas-mvp/packages/dal/types"

type Model struct {
	TableName    string
	query        string
	ReturnValue  []types.DALType
	Error        error
	LastInsertID int64
}

func (m *Model) Query() string {
	m.query += ";"
	return m.query
}
