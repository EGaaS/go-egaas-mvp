package types

type DALByteArray struct {
	dataSource DataSource
	value      string
	isNull     bool
}

func (t DALByteArray) DataSource() DataSource {
	return t.dataSource
}

func (t DALByteArray) Value() string {
	return t.value
}

func (t DALByteArray) IsNull() bool {
	return t.isNull
}

func (t DALByteArray) SetNull(null bool) {
	t.isNull = null
}

func (t DALByteArray) Set(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = string(data)
	}
}

func (t DALByteArray) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = string(data)
	}
}

func (t DALByteArray) String() string {
	return string(t.value)
}
