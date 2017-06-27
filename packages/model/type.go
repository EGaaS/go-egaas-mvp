package model

import "github.com/shopspring/decimal"

import "fmt"
import "encoding/binary"

type OperationType byte

const (
	Create OperationType = iota
	Read
	Update
	Delete
)

type Operation struct {
	Type    OperationType
	ColName string
}

type DbType interface {
	ColName() string
	String() string
	IsNull() bool
	SetNull(bool)
	FromBytes([]byte)
}

type dbInt64 struct {
	value   int64
	colName string
	isNull  bool
}

type dbInt32 struct {
	value   int32
	colName string
	isNull  bool
}

type dbString struct {
	value   string
	colName string
	isNull  bool
}

type dbDecimal struct {
	value   decimal.Decimal
	colName string
	isNull  bool
}

// int64
func (dbi dbInt64) ColName() string {
	return dbi.colName
}

func (dbi dbInt64) Value() int64 {
	return dbi.value
}

func (dbi dbInt64) String() string {
	return fmt.Sprintf("%d", dbi.value)
}

func (dbi dbInt64) IsNull() bool {
	return dbi.isNull
}

func (dbi dbInt64) SetNull(null bool) {
	dbi.isNull = null
}

func (dbi dbInt64) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		dbi.isNull = true
	} else {
		dbi.isNull = false
		dbi.value = int64(binary.BigEndian.Uint64(data))
	}
}

// dbDecimal
func (dbd dbDecimal) ColName() string {
	return dbd.colName
}

func (dbd dbDecimal) Value() decimal.Decimal {
	return dbd.value
}

func (dbd dbDecimal) SetValue(val decimal.Decimal) dbDecimal {
	dbd.value = val
	return dbd
}

func (dbd dbDecimal) String() string {
	return dbd.value.String()
}

func (dbd dbDecimal) IsNull() bool {
	return dbd.isNull
}

func (dbd dbDecimal) SetNull(null bool) {
	dbd.isNull = null
}

//TODO доделать конвертацию
func (dbd dbDecimal) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		dbd.isNull = true
	} else {
		dbd.isNull = false
	}
}

//dbInt32
func (dbi dbInt32) ColName() string {
	return dbi.colName
}

func (dbi dbInt32) Value() int32 {
	return dbi.value
}

func (dbi dbInt32) String() string {
	return fmt.Sprintf("%d", dbi.value)
}

func (dbi dbInt32) IsNull() bool {
	return dbi.isNull
}

func (dbi dbInt32) SetNull(null bool) {
	dbi.isNull = null
}

func (dbi dbInt32) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		dbi.isNull = true
	} else {
		dbi.isNull = false
		dbi.value = int32(binary.BigEndian.Uint32(data))
	}
}

//dbString
func (dbs dbString) ColName() string {
	return dbs.colName
}

func (dbs dbString) Value() string {
	return dbs.value
}

func (dbs dbString) String() string {
	return dbs.value
}

func (dbs dbString) IsNull() bool {
	return dbs.isNull
}

func (dbs dbString) SetNull(null bool) {
	dbs.isNull = null
}

func (dbs dbString) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		dbs.isNull = true
	} else {
		dbs.isNull = false
		dbs.value = string(data)
	}
}
