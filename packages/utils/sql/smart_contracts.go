package sql

func (db *DCDB) GetOrderedSmartContracts(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_smart_contracts" order by id`, -1)
}

func (db *DCDB) GetSmartContractsByID(tablePrefix string, contractID int64) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_smart_contracts" WHERE id = ?`, contractID).String()
}

func (db *DCDB) GetSmartContractsByName(tablePrefix string, contractName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_smart_contracts" WHERE name = ?`, contractName).String()
}
