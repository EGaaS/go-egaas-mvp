package types

type DALType interface {
	SetDatasource(DataSource)
	DataSource() DataSource
	IsNull() bool
	SetNull(bool)
	FromBytes([]byte)
	String() string
}
