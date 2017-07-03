package types

type DALByteArray struct {
	dataSource DataSource
	isNull     bool
	value      []byte
}

func (t DALByteArray) Set(data []byte) DALByteArray {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = data
	}

	return t
}

func (t DALByteArray) SetDatasource(ds DataSource) DALByteArray {
	t.dataSource = ds
	return t
}

func (t DALByteArray) SetNull(null bool) DALByteArray {
	t.isNull = null
	return t
}

func (t DALByteArray) FromBytes(data []byte) DALByteArray {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = data
	}
	return t
}

func (t DALByteArray) IsNull() bool {
	return t.isNull
}

func (t DALByteArray) Value() []byte {
	return t.value
}

func (t DALByteArray) String() string {
	return string(t.value)
}

func (t DALByteArray) DataSource() DataSource {
	return t.dataSource
}
