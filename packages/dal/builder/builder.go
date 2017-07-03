package builder

import (
	"github.com/EGaaS/go-egaas-mvp/packages/dal/model"
	"github.com/EGaaS/go-egaas-mvp/packages/dal/types"
)

type Builder interface {
	Create(model.Model)
	Read(model.Model)
	Update(model.Model)
	Delete(model.Model)
	Where(Condition)
	And(Condition)
	Or(Condition)
	Compile() string
}

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
