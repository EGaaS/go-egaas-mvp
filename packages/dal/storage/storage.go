package storage

import (
	"github.com/EGaaS/go-egaas-mvp/packages/dal/model"
	"github.com/EGaaS/go-egaas-mvp/packages/dal/types"
)

type DataProvider byte

const (
	POSTGRES DataProvider = iota
)

type Comparator byte

const (
	Less Comparator = iota
	LessOrEqual
	Equal
	NotEqual
	GreaterOrEqual
	Greater
)

type Condition struct {
	Field      types.DALType
	Comparator Comparator
	Value      string
}

func NewStorage(prov DataProvider) *Storage {
	switch prov {
	case POSTGRES:
		if result, err := PgConnect("login", "pass", "test", 5432); err != nil {
			return nil
		}
		return result
	}
	return nil
}

type Storage interface {
	Connect()
	Create(model *model.Model) error
	Read()
	Update()
	Delete()
}

type Builder interface {
	Where()
	And()
	Or()
}
