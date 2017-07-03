package types

import (
	"encoding/binary"
	"fmt"
)

type DALInt64 struct {
	dataSource DataSource
	isNull     bool
	value      int64
}

func (t *DALInt64) SetDatasource(ds DataSource) {
	t.dataSource = ds
}

func (t *DALInt64) Set(val int64) {
	t.value = val
}

func (t *DALInt64) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
		t.value = int64(binary.BigEndian.Uint64(data))
	}
	return
}

func (t *DALInt64) SetNull(null bool) {
	t.isNull = null
}

func (t *DALInt64) Value() int64 {
	return t.value
}

func (t *DALInt64) String() string {
	return fmt.Sprintf("%d", t.value)
}

func (t *DALInt64) DataSource() DataSource {
	return t.dataSource
}

func (t *DALInt64) IsNull() bool {
	return t.isNull
}
