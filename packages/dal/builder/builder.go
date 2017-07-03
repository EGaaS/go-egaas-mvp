package builder

import (
	"github.com/EGaaS/go-egaas-mvp/packages/dal/model"
	"github.com/EGaaS/go-egaas-mvp/packages/dal/storage"
)

type Builder interface {
	Create(model.Model)
	Read(model.Model)
	Update(model.Model)
	Delete(model.Model)
	Where(storage.Condition)
	And(storage.Condition)
	Or(storage.Condition)
	Compile() string
}
