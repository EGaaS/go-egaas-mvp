package sql

func (db *DCDB) GetDltTransferPrice() (int64, error) {
	return db.Single(`SELECT value->'dlt_transfer' FROM system_parameters WHERE name = ?`, "op_price").Int64()
}

func (db *DCDB) GetDltTransferPriceString() (string, error) {
	return db.Single(`SELECT value->'dlt_transfer' FROM system_parameters WHERE name = ?`, "op_price").String()
}

func (db *DCDB) GetAllSystemParameters() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM system_parameters`, -1)
}
