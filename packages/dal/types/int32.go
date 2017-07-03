package types

import (
	"encoding/binary"
	"fmt"
)

type DALInt32 struct {
	dataSource DataSource
	isNull     bool
	value      int32
}

func (t DALInt32) SetDatasource(ds DataSource) DALInt32 {
	t.dataSource = ds
	return t
}

func (t DALInt32) SetNull(null bool) DALInt32 {
	t.isNull = null
	return t
}

func (t DALInt32) Set(val int32) DALInt32 {
	t.value = val
	return t
}

func (t DALInt32) FromBytes(data []byte) DALInt32 {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = int32(binary.BigEndian.Uint32(data))
	}
	return t
}

func (t DALInt32) DataSource() DataSource {
	return t.dataSource
}

func (t DALInt32) IsNull() bool {
	return t.isNull
}

func (t DALInt32) Value() int32 {
	return t.value
}

func (t DALInt32) String() string {
	return fmt.Sprintf("%d", t.value)
}
