package builder

import (
	"strings"

	"github.com/EGaaS/go-egaas-mvp/packages/dal/types"
)

type PgBuilder struct {
	Querys []string
}

func (pg *PgBuilder) splitTables(fields []types.DALType) map[string][]types.DALType {
	tables := make(map[string][]types.DALType, 0)

	for _, field := range fields {
		inner_tables := strings.Split(field.DataSource().ResName, "|")
		for _, inner_table := range inner_tables {
			if tables[inner_table] == nil {
				tables[inner_table] = []types.DALType{field}
			} else {
				tables[inner_table] = append(tables[field.DataSource().ResName], field)
			}
		}
	}
	return tables
}

func (pg *PgBuilder) Create(fields ...types.DALType) *PgBuilder {
	tables := pg.splitTables(fields)

	for table, values := range tables {
		query := "insert into " + table
		columns := " ("
		vals := "values ("
		for _, value := range values {
			columns += value.DataSource().ParamName + ", "
			vals += value.String() + ", "
		}
		query += columns[:len(columns)-2] + ") " + vals[:len(vals)-2] + ")"
		pg.Querys = append(pg.Querys, query)
	}
	return pg
}

func (pg *PgBuilder) Read(fields ...types.DALType) *PgBuilder {
	tables := pg.splitTables(fields)

	for table, values := range tables {
		query := "select "
		for _, value := range values {
			resources := strings.Split(value.DataSource().ResName, "|")
			if len(resources) == 1 {
				query += value.DataSource().ParamName + ", "
			}
			if len(resources) > 1 && resources[0] == table {
				query += value.DataSource().ParamName + ", "
			}
		}
		query = query[:len(query)-2] + " from " + table
		pg.Querys = append(pg.Querys, query)
	}
	return pg
}

func (pg *PgBuilder) Update(fields ...types.DALType) *PgBuilder {
	tables := pg.splitTables(fields)

	for table, values := range tables {
		query := "update " + table + " set "
		for _, value := range values {
			resources := strings.Split(value.DataSource().ResName, "|")
			if len(resources) == 1 {
				query += value.DataSource().ParamName + " = " + value.String() + ", "
			}
			if len(resources) > 1 && resources[0] == table {
				query += value.DataSource().ParamName + " = " + value.String() + ", "
			}
		}
		query = query[:len(query)-2]
		pg.Querys = append(pg.Querys, query)
	}
	return pg
}

func (pg *PgBuilder) Delete(fields ...types.DALType) *PgBuilder {
	tables := pg.splitTables(fields)
	for table := range tables {
		query := "delete from " + table
		pg.Querys = append(pg.Querys, query)
	}
	return pg
}

func (pg *PgBuilder) Where(condition Condition) *PgBuilder {
	for i := 0; i < len(pg.Querys); i++ {
		if strings.Contains(pg.Querys[i], " "+condition.Field.DataSource().ParamName) {
			pg.Querys[i] += ";"
		} else {
			pg.Querys[i] += " where "
			switch condition.Comparator {
			case Less:
				pg.Querys[i] += condition.Field.DataSource().ParamName + " < " + condition.Value
			case LessOrEqual:
				pg.Querys[i] += condition.Field.DataSource().ParamName + " <= " + condition.Value
			case Equal:
				pg.Querys[i] += condition.Field.DataSource().ParamName + " = " + condition.Value
			case NotEqual:
				pg.Querys[i] += condition.Field.DataSource().ParamName + " != " + condition.Value
			case GreaterOrEqual:
				pg.Querys[i] += condition.Field.DataSource().ParamName + " >= " + condition.Value
			case Greater:
				pg.Querys[i] += condition.Field.DataSource().ResName + " > " + condition.Value
			}
		}
	}
	return pg
}

func (pg *PgBuilder) And(condition Condition) *PgBuilder {
	for i := 0; i < len(pg.Querys); i++ {
		if pg.Querys[i][len(pg.Querys[i])-1:] == ";" {
			continue
		}

		pg.Querys[i] += " and "
		switch condition.Comparator {
		case Less:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " < " + condition.Value
		case LessOrEqual:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " <= " + condition.Value
		case Equal:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " - " + condition.Value
		case NotEqual:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " != " + condition.Value
		case GreaterOrEqual:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " >= " + condition.Value
		case Greater:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " > " + condition.Value
		}
	}
	return pg
}

func (pg *PgBuilder) Or(condition Condition) *PgBuilder {
	for i := 0; i < len(pg.Querys); i++ {
		if pg.Querys[i][len(pg.Querys[i])-1:] == ";" {
			continue
		}

		pg.Querys[i] += " or "
		switch condition.Comparator {
		case Less:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " < " + condition.Value
		case LessOrEqual:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " <= " + condition.Value
		case Equal:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " - " + condition.Value
		case NotEqual:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " != " + condition.Value
		case GreaterOrEqual:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " >= " + condition.Value
		case Greater:
			pg.Querys[i] += condition.Field.DataSource().ParamName + " > " + condition.Value
		}
	}
	return pg
}

func (pg *PgBuilder) Compile() *PgBuilder {
	for i := 0; i < len(pg.Querys); i++ {
		if pg.Querys[i][len(pg.Querys[i])-1:] != ";" {
			pg.Querys[i] += ";"
		}
	}
	return pg
}
