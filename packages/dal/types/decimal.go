package types

import "github.com/shopspring/decimal"

type DALDecimal struct {
	dataSource DataSource
	isNull     bool
	value      decimal.Decimal
}

func (t DALDecimal) SetDatasource(ds DataSource) DALDecimal {
	t.dataSource = ds
	return t
}

func (t DALDecimal) SetNull(null bool) DALDecimal {
	t.isNull = null
	return t
}

//TODO доделать конвертацию
func (t DALDecimal) FromBytes(data []byte) DALDecimal {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
	}
	return t
}

func (t DALDecimal) Set(val decimal.Decimal) DALDecimal {
	t.value = val
	return t
}

func (t DALDecimal) DataSource() DataSource {
	return t.dataSource
}

func (t DALDecimal) IsNull() bool {
	return t.isNull
}

func (t DALDecimal) Value() decimal.Decimal {
	return t.value
}

func (t DALDecimal) String() string {
	return t.value.String()
}
