package types

type DALString struct {
	dataSource DataSource
	isNull     bool
	value      string
}

func (t DALString) SetDatasource(ds DataSource) DALString {
	t.dataSource = ds
	return t
}

func (t DALString) Set(val string) DALString {
	t.value = val
	return t
}

func (t DALString) FromBytes(data []byte) DALString {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = string(data)
	}
	return t
}

func (t DALString) SetNull(null bool) DALString {
	t.isNull = null
	return t
}

func (t DALString) DataSource() DataSource {
	return t.dataSource
}

func (t DALString) IsNull() bool {
	return t.isNull
}

func (t DALString) Value() string {
	return t.value
}

func (t DALString) String() string {
	return t.value
}
