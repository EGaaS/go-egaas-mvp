package postgres

import (
	"github.com/EGaaS/go-egaas-mvp/packages/dal/model"
)

func (m *model.Model) Where(condition Condition) *model.Model {
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

func (m *model.Model) And(condition Condition) *model.Model {
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

func (m *model.Model) Or(condition Condition) *model.Model {
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
