package types

// DataSource All DAL types knows where it could be stored or retrieved.
// ResName is the address of the resource unit which contains the value,
// for example table name in SQL database or URL for http request
// ParamName is the address of concrete point where the value is stored,
// like a column in SQL table or a parameter of the http request
type DataSource struct {
	ResName   string
	ParamName string
}
