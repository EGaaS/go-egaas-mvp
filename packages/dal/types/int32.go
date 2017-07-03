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

func (t *DALInt32) SetDatasource(ds DataSource) {
	t.dataSource = ds
}

func (t *DALInt32) SetNull(null bool) {
	t.isNull = null
}

func (t *DALInt32) Set(val int32) {
	t.value = val
}

func (t *DALInt32) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = int32(binary.BigEndian.Uint32(data))
	}
}

func (t *DALInt32) DataSource() DataSource {
	return t.dataSource
}

func (t *DALInt32) IsNull() bool {
	return t.isNull
}

func (t *DALInt32) Value() int32 {
	return t.value
}

func (t *DALInt32) String() string {
	return fmt.Sprintf("%d", t.value)
}
