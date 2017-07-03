package types

import "github.com/shopspring/decimal"

type DALDecimal struct {
	dataSource DataSource
	isNull     bool
	value      decimal.Decimal
}

func (t *DALDecimal) SetDatasource(ds DataSource) {
	t.dataSource = ds
}

func (t *DALDecimal) SetNull(null bool) {
	t.isNull = null
}

//TODO доделать конвертацию
func (t *DALDecimal) FromBytes(data []byte) {
	if data == nil || len(data) == 0 {
		t.isNull = true
	} else {
		t.isNull = false
	}
}

func (t *DALDecimal) Set(val decimal.Decimal) {
	t.value = val
}

func (t *DALDecimal) DataSource() DataSource {
	return t.dataSource
}

func (t *DALDecimal) IsNull() bool {
	return t.isNull
}

func (t *DALDecimal) Value() decimal.Decimal {
	return t.value
}

func (t *DALDecimal) String() string {
	return t.value.String()
}
