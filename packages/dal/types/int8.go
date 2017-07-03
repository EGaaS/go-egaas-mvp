package types

import (
	"encoding/binary"
	"fmt"
)

type DALInt8 struct {
	dataSource DataSource
	isNull     bool
	value      int8
}

func (t *DALInt8) SetDatasource(ds DataSource) {
	t.dataSource = ds
}
func (t *DALInt8) SetNull(null bool) {
	t.isNull = null
}

func (t *DALInt8) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = int8(binary.BigEndian.Uint16(data))
	}
}

func (t *DALInt8) Set(val int8) {
	t.value = val
}

func (t *DALInt8) DataSource() DataSource {
	return t.dataSource
}

func (t *DALInt8) IsNull() bool {
	return t.isNull
}

func (t *DALInt8) Value() int8 {
	return t.value
}

func (t *DALInt8) String() string {
	return fmt.Sprintf("%d", t.value)
}
