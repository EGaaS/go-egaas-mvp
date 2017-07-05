package sql

import "fmt"

func (db *DCDB) GetAppsByName(tableName, name string) (map[string]string, error) {
	return db.OneRow(`select * from `+tableName+` where name=?`, name).String()
}

func (db *DCDB) UpdateBlockInfoInApps(tableName string, blockID int64, done int, name string) error {
	return db.ExecSQL(fmt.Sprintf(`update %s set done=?, blocks=concat(blocks, ',%d') where name=?`, tableName, blockID),
		done, name)
}

func (db *DCDB) GetNameDoneFromApps(tableName string) ([]map[string]string, error) {
	return db.GetAll(`select name,done from `+tableName, -1)
}

func (db *DCDB) CreateBlockInfoInApps(tableName string, blockID int64, done int, name string) error {
	return db.ExecSQL(fmt.Sprintf(`insert into %s (name,done,blocks) values(?,?,'%d')`, tableName, blockID),
		name, done)
}

func (db *DCDB) GetTablePermissions(tablePrefix string, tableName string) (map[string]string, error) {
	return db.GetMap(`SELECT data.* FROM "`+tablePrefix+`_tables", jsonb_each_text(columns_and_permissions) as data WHERE name = ?`, "key", "value", tableName)
}

func (db *DCDB) GetColumnsAndPermissions(tablePrefix string, tableName string) (map[string]string, error) {
	return db.GetMap(`SELECT data.* FROM "`+tablePrefix+`_tables", jsonb_each_text(columns_and_permissions->'update') as data WHERE name = ?`, "key", "value", tableName)
}

func (db *DCDB) GetColumnPermissionsByState(stateName string, tableName string) (string, error) {
	return db.Single(fmt.Sprintf(`select columns_and_permissions->'update' from "%s_tables" where name=?`, stateName), tableName).String()
}

func (db *DCDB) GetPermissionsByState(stateName string, tableName string) (string, error) {
	return db.Single(fmt.Sprintf(`select columns_and_permissions from "%s_tables" where name=?`, stateName), tableName).String()
}

func (db *DCDB) GetFullTableData(tablePrefix string, tableName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_tables" WHERE name = ?`, tableName).String()
}

func (db *DCDB) Get1000RecFromTable(tableName string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tableName+`" order by id`, 1000)
}

func (db *DCDB) GetContractConditionsAndValue(stateName, contractName, contract string) (map[string]string, error) {
	return db.OneRow(fmt.Sprintf(`select conditions,value from "%s_%s" where name=?`, stateName, contractName), contract).String()
}

func (db *DCDB) GetTables(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_tables"`, -1)
}

func (db *DCDB) GetRecordsCount(tableName string) (int64, error) {
	return db.Single(`SELECT count(id) FROM "` + tableName + `"`).Int64()
}
