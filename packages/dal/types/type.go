package types

type DALType interface {
	DataSource() DataSource
	IsNull() bool
	SetNull(bool)
	FromBytes([]byte)
	String() string
}
