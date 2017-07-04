package sql

import "fmt"

func (db *DCDB) GetAppsByName(tableName, name string) (map[string]string, error) {
	return db.OneRow(`select * from `+tableName+` where name=?`, name).String()
}

func (db *DCDB) UpdateBlockInfoInApps(tableName string, blockID int64, done int, name string) error {
	return db.ExecSQL(fmt.Sprintf(`update %s set done=?, blocks=concat(blocks, ',%d') where name=?`, tableName, blockID),
		done, name)
}

func (db *DCDB) CreateBlockInfoInApps(tableName string, blockID int64, done int, name string) error {
	return db.ExecSQL(fmt.Sprintf(`insert into %s (name,done,blocks) values(?,?,'%d')`, tableName, blockID),
		name, done)
}

func (db *DCDB) GetDltTransfer() (int64, error) {
	return db.Single(`SELECT value->'dlt_transfer' FROM system_parameters WHERE name = ?`, "op_price").Int64()
}

func (db *DCDB) GetNameDoneFromApps(tableName string) ([]map[string]string, error) {
	return db.GetAll(`select name,done from `+tableName, -1)
}

func (db *DCDB) GetInstallationState() (string, error) {
	return db.Single("SELECT progress FROM install").String()
}

func (db *DCDB) GetFirstLoadBlockchainURL() (string, error) {
	return db.Single("SELECT first_load_blockchain_url FROM config").String()
}

func (db *DCDB) GetTablePermissions(tablePrefix string, tableName string) (map[string]string, error) {
	return db.GetMap(`SELECT data.* FROM "`+tablePrefix+`_tables", jsonb_each_text(columns_and_permissions) as data WHERE name = ?`, "key", "value", tableName)
}

func (db *DCDB) GetColumnsAndPermissions(tablePrefix string, tableName string) (map[string]string, error) {
	return db.GetMap(`SELECT data.* FROM "`+tablePrefix+`_tables", jsonb_each_text(columns_and_permissions->'update') as data WHERE name = ?`, "key", "value", tableName)
}

func (db *DCDB) GetFullTableData(tablePrefix string, tableName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_tables" WHERE name = ?`, tableName).String()
}

func (db *DCDB) GetColumnsCount(tableName string) (int64, error) {
	return db.Single("SELECT count(column_name) FROM information_schema.columns WHERE table_name=?", tableName).Int64()
}
