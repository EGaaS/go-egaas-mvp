package types

type DALByteArray struct {
	dataSource DataSource
	isNull     bool
	value      []byte
}

func (t DALByteArray) Set(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = data
	}
}

func (t *DALByteArray) SetDatasource(ds DataSource) {
	t.dataSource = ds
}

func (t *DALByteArray) SetNull(null bool) {
	t.isNull = null
}

func (t *DALByteArray) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = data
	}
}

func (t *DALByteArray) IsNull() bool {
	return t.isNull
}

func (t *DALByteArray) Value() []byte {
	return t.value
}

func (t *DALByteArray) String() string {
	return string(t.value)
}

func (t *DALByteArray) DataSource() DataSource {
	return t.dataSource
}
