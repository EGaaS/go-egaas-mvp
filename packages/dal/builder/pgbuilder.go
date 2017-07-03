package builder

import (
	"reflect"

	"github.com/EGaaS/go-egaas-mvp/packages/dal/model"
	"strings"
)

type PgBuilder struct {
	query []string
}

func (pg *PgBuilder) Create(m model.Model) error {
	q := "insert into " + m.TableName
	columns := " ("
	values := "values ("
	for _, field := range fields {
		columns += field.ColName() + ", "
		values += field.String() + ", "
	}
	q += columns[:len(columns)-2] + ") " + values[:len(values)-2] + ")"
	pg.query = query
	return nil
}

func (pg *PgBuilder) collectModels(m model.Model) []model.Model {
	v := reflect.ValueOf(m)

	for i := 0; i < v.NumField(); i++ {
		if strings.Contains(v.Field(i).Type(), "model") {
			
		}
	}
}

func (pg *PgBuilder) Read(m model.Model) *PgBuilder {
	q := "select "
	for _, field := range fields {
		m.ReturnValue = append(m.ReturnValue, field)
		q += field.ColName() + ", "
	}
	q = q[:len(query)-2] + " from " + m.TableName
	pg.query = query
	return pg
}

func (pg *PgBuilder) Update(m model.Model) PgBuilder {

}

func (pg *PgBuilder) Delete(m model.Model) PgBuilder {
	query := "delete from " + m.TableName + " where "
	for _, condition := range conditions {
		query += condition.Field.ColName()
		switch condition.Comparator {
		case Less:
			query += " < "
		case LessOrEqual:
			query += " > "
		case Equal:
			query += " = "
		case GreaterOrEqual:
			query += " >= "
		case Greater:
			query += " > "
		}
		query += condition.Value + " and "
	}
	m.query += query[:len(query)-5]
	return m
}

func (pg *PgBuilder) Where(condition Condition) PgBuilder {
	query := " where "
	switch condition.Comparator {
	case Less:
		query += condition.Field.ColName() + " < " + condition.Value
	case LessOrEqual:
		query += condition.Field.ColName() + " <= " + condition.Value
	case Equal:
		query += condition.Field.ColName() + " = " + condition.Value
	case NotEqual:
		query += condition.Field.ColName() + " != " + condition.Value
	case GreaterOrEqual:
		query += condition.Field.ColName() + " >= " + condition.Value
	case Greater:
		query += condition.Field.ColName() + " > " + condition.Value
	}
	m.query += query
	return m
}

func (pg *PgBuilder) And(condition Condition) *PgBuilder {
	query := " and "
	switch condition.Comparator {
	case Less:
		query += condition.Field.ColName() + " < " + condition.Value
	case LessOrEqual:
		query += condition.Field.ColName() + " <= " + condition.Value
	case Equal:
		query += condition.Field.ColName() + " - " + condition.Value
	case NotEqual:
		query += condition.Field.ColName() + " != " + condition.Value
	case GreaterOrEqual:
		query += condition.Field.ColName() + " >= " + condition.Value
	case Greater:
		query += condition.Field.ColName() + " > " + condition.Value
	}
	m.query += query
	return m
}

func (pg *PgBuilder) Or(condition Condition) *PgBuilder {
	query := " or "
	switch condition.Comparator {
	case Less:
		query += condition.Field.ColName() + " < " + condition.Value
	case LessOrEqual:
		query += condition.Field.ColName() + " <= " + condition.Value
	case Equal:
		query += condition.Field.ColName() + " - " + condition.Value
	case NotEqual:
		query += condition.Field.ColName() + " != " + condition.Value
	case GreaterOrEqual:
		query += condition.Field.ColName() + " >= " + condition.Value
	case Greater:
		query += condition.Field.ColName() + " > " + condition.Value
	}
	m.query += query
	return m
}

func (pg *PgWorker) Compile() []string {

}
