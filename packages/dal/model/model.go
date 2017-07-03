package model

import (
	"errors"

	"github.com/EGaaS/go-egaas-mvp/packages/dal/types"
)

var (
	UnknownDatasorceField = errors.New("Unknown field in selected database")
)

type Model struct {
	ID           types.DALString
	LastInsertID types.DALString
}
