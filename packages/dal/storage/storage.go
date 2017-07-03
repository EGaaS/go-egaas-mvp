package storage

import "github.com/EGaaS/go-egaas-mvp/packages/dal/model"

type DataProvider byte

const (
	POSTGRES DataProvider = iota
)

/*
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
*/
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
