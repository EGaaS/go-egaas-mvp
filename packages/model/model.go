package model

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
	Field      DbType
	Comparator Comparator
	Value      string
}

type Model struct {
	TableName    string
	query        string
	ReturnValue  []DbType
	Error        error
	LastInsertID int64
}

func (m *Model) Read(fields ...DbType) *Model {
	query := "select "
	for _, field := range fields {
		m.ReturnValue = append(m.ReturnValue, field)
		query += field.ColName() + ", "
	}
	query = query[:len(query)-2] + " from " + m.TableName
	m.query = query
	return m
}

func (m *Model) Create(fields ...DbType) *Model {
	query := "insert into " + m.TableName
	columns := " ("
	values := "values ("
	for _, field := range fields {
		columns += field.ColName() + ", "
		values += field.String() + ", "
	}
	query += columns[:len(columns)-2] + ") " + values[:len(values)-2] + ")"
	m.query = query
	return m
}

func (m *Model) Where(condition Condition) *Model {
	query := " where "
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

func (m *Model) And(condition Condition) *Model {
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

func (m *Model) Or(condition Condition) *Model {
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

func (m *Model) Delete(conditions ...Condition) *Model {
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

func (m *Model) Query() string {
	m.query += ";"
	return m.query
}
