package sql

func (db *DCDB) GetDltTransfer() (int64, error) {
	return db.Single(`SELECT value->'dlt_transfer' FROM system_parameters WHERE name = ?`, "op_price").Int64()
}

func (db *DCDB) GetAllSystemParameters() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM system_parameters`, -1)
}
