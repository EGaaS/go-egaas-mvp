package types

type DALString struct {
	dataSource DataSource
	value      string
	isNull     bool
}

func (t DALString) DataSource() DataSource {
	return t.dataSource
}

func (t DALString) Value() string {
	return t.value
}

func (t DALString) String() string {
	return t.value
}

func (t DALString) IsNull() bool {
	return t.isNull
}

func (t DALString) SetNull(null bool) {
	t.isNull = null
}

func (t DALString) Set(val string) {
	t.value = val
}

func (t DALString) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = string(data)
	}
}
