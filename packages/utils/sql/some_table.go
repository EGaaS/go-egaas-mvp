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

func (db *DCDB) GetConditionsFromSomewhere(tablePrefix string, tableName string, conditionName string) (string, error) {
	return db.Single(`SELECT conditions FROM "`+tablePrefix+`_`+tableName+`" WHERE name = ?`, conditionName).String()
}

func (db *DCDB) GetTables(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_tables"`, -1)
}

func (db *DCDB) GetRecordsCount(tableName string) (int64, error) {
	return db.Single(`SELECT count(id) FROM "` + tableName + `"`).Int64()
}

func (db *DCDB) GetAnotherRecordsCount(tableName string) (int64, error) {
	return db.Single(`SELECT count(*) FROM ` + tableName).Int64()
}

func (db *DCDB) GetRecordsCountWhereName(tableName string, anotherTableName string) (int64, error) {
	return db.Single(`select count(*) from "`+tableName+`" where name = ?`, anotherTableName).Int64()
}

func (db *DCDB) GetRecordsFromSomwhereWithSomethingID(columnName string, tableName string, idName string, id string) (string, error) {
	return db.Single(`select `+columnName+` from `+tableName+` where `+idName+`=?`, id).String()
}

func (db *DCDB) GetCustomFieldsFromCustomTable(fields string, tableNameAndClause string) ([]map[string]string, error) {
	return db.GetAll(`select `+fields+` from `+tableNameAndClause, -1)
}

func (db *DCDB) GetAmountFromSomewhere(tableName string, columnName string, id int64) (string, error) {
	return db.Single("SELECT amount FROM "+tableName+" WHERE "+columnName+" = ?", id).String()
}

func (db *DCDB) GetCustomFieldsFromCustomTableWithLimit(fields string, tableNameAndClause string, limit int) ([]map[string]string, error) {
	return db.GetAll(`select `+fields+` from `+tableNameAndClause, limit)
}

func (db *DCDB) GetCustomFieldsFromCustomTableOneRow(fields string, tableNameAndClause string) (map[string]string, error) {
	return db.OneRow(`select ` + fields + ` from ` + tableNameAndClause).String()
}

func (db *DCDB) GetCustomFieldsFromCustomTableOneField(field string, tableNameAndClause string) (string, error) {
	return db.Single(`select ` + field + ` from ` + tableNameAndClause).String()
}

func (db *DCDB) DeleteAllFromTable(tableName string) error {
	return db.ExecSQL(`DELETE FROM ` + tableName)
}

func (db *DCDB) SelectMaxAutoincrementID(columnName string, tableName string) (int64, error) {
	return db.Single("SELECT " + columnName + " FROM " + tableName + " ORDER BY " + columnName + " DESC LIMIT 1").Int64()
}

func (db *DCDB) GetPgSerialSequence(columnName string, tableName string) (string, error) {
	return db.Single("SELECT pg_get_serial_sequence('" + tableName + "', '" + columnName + "')").String()
}

func (db *DCDB) PgRestartSequence(sequence string, newAutoincrement string) error {
	return db.ExecSQL("ALTER SEQUENCE " + sequence + " RESTART WITH " + newAutoincrement)
}

func (db *DCDB) MySQLRestartSequence(tableName string, newAutoincrement string) error {
	return db.ExecSQL("ALTER TABLE " + tableName + " AUTO_INCREMENT = " + newAutoincrement)
}

func (db *DCDB) SQLiteRestartSequence(tableName string, newAutoincrement int64) error {
	return db.ExecSQL("UPDATE SQLITE_SEQUENCE SET seq = ? WHERE name = ?", newAutoincrement, tableName)
}

func (db *DCDB) UpdateSomeTable(tableName string, values string, where string, rbID string) error {
	return db.ExecSQL(`UPDATE "`+tableName+`" SET `+values+` rb_id = ? `+where, rbID)
}

func (db *DCDB) SelectRbIDFromSomeTable(tableName string, whereClause string) (int64, error) {
	return db.Single("SELECT rb_id FROM " + tableName + " " + whereClause + " order by rb_id desc").Int64()
}

func (db *DCDB) DeleteSomethingFromSomewhere(tableName string, whereClause string) error {
	return db.ExecSQL("DELETE FROM " + tableName + " " + whereClause)
}

func (db *DCDB) DeleteFromTable(tablePrefix string, tableName string) error {
	return db.ExecSQL(`DELETE FROM "`+tablePrefix+`_tables" WHERE name = ?`, tableName)
}

func (db *DCDB) IsSqliteTableExists(tableName string) (int64, error) {
	return db.Single(`SELECT name FROM sqlite_master WHERE type='table' AND name='` + tableName + `';`).Int64()
}

func (db *DCDB) IsMySQLTableExists(tableName string) (int64, error) {
	return db.Single(`SHOW TABLES LIKE '` + tableName + `'`).Int64()
}

func (db *DCDB) IsPGTableExists(tableName string) (int64, error) {
	return db.Single(`SELECT relname FROM pg_class WHERE relname = '` + tableName + `';`).Int64()
}

func (db *DCDB) IsColumnExists(tableName string, columnName string, anotherTableName string) (int64, error) {
	return db.Single(`select count(*) from "`+tableName+`" where (columns_and_permissions->'update'-> ? ) is not null AND name = ?`, columnName, anotherTableName).Int64()
}

func (db *DCDB) GetPermissionsAndRollbackID(tableName string, columnName string, anotherTableName string) (map[string]string, error) {
	return db.OneRow(`SELECT columns_and_permissions, rb_id FROM "`+tableName+`" where (columns_and_permissions->'update'-> ? ) is not null AND name = ?`,
		columnName, anotherTableName).String()
}

func (db *DCDB) GetPermissionsAndRollbackIDByName(tableName string, anotherTableName string) (map[string]string, error) {
	return db.OneRow(`SELECT columns_and_permissions, rb_id FROM "`+tableName+`" where name=?`, anotherTableName).String()
}
func (db *DCDB) UpdateColumnAndPermissions(tableName string, columnName string, permissions string, rollbackID string, anotherTableName string) error {
	return db.ExecSQL(`UPDATE "`+tableName+`" SET columns_and_permissions = jsonb_set(columns_and_permissions, '{update, `+columnName+`}', ?, true), rb_id = ? WHERE name = ?`,
		`"`+permissions+`"`, rollbackID, anotherTableName)
}

func (db *DCDB) CreateColumnsAndPermissions(tableName string, actionName string, actionValue string, rollbackID string, anotherTableName string) error {
	return db.ExecSQL(`UPDATE "`+tableName+`" SET columns_and_permissions = jsonb_set(columns_and_permissions, '{`+actionName+`}', ?, true), rb_id = ? WHERE name = ?`,
		`"`+actionValue+`"`, rollbackID, anotherTableName)
}

func (db *DCDB) GetSomethingSomwhereByID(columnName string, tableName string, id int64) (int64, error) {
	return db.Single(`select `+columnName+` from `+tableName+` where id=?`, id).Int64()
}

func (db *DCDB) GetSomethingFromSomewhereSomehow(columnName string, tableName string, whereClause string, limit int, params ...interface{}) ([]map[string]string, error) {
	return db.GetAll(`select `+columnName+` from `+tableName+` where `+
		whereClause, int(limit), params...)
}

func (db *DCDB) CreateTables(id string) error {
	return db.ExecSQL(`CREATE TABLE "` + id + `_tables" (
				"name" varchar(100)  NOT NULL DEFAULT '',
				"columns_and_permissions" jsonb,
				"conditions" text  NOT NULL DEFAULT '',
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER TABLE ONLY "` + id + `_tables" ADD CONSTRAINT "` + id + `_tables_pkey" PRIMARY KEY (name);
				`)
}

func (db *DCDB) CreateTablesRecords(id string, sid string, psid string) error {
	return db.ExecSQL(`INSERT INTO "`+id+`_tables" (name, columns_and_permissions, conditions) VALUES
		(?, ?, ?)`,
		id+`_citizens`, `{"general_update":"`+sid+`", "update": {"public_key_0": "`+sid+`"}, "insert": "`+sid+`", "new_column":"`+sid+`"}`, psid)
}

func (db *DCDB) CreateAppsTable(id string) error {
	return db.ExecSQL(`CREATE TABLE "` + id + `_apps" (
				"name" varchar(100)  NOT NULL DEFAULT '',
				"done" integer NOT NULL DEFAULT '0',
				"blocks" text  NOT NULL DEFAULT ''
				);
				ALTER TABLE ONLY "` + id + `_apps" ADD CONSTRAINT "` + id + `_apps_pkey" PRIMARY KEY (name);
				`)
}

func (db *DCDB) DropTable(tableID int64, tableName string) error {
	return db.ExecSQL(fmt.Sprintf(`DROP TABLE "%d_%s"`, tableID, tableName))
}

func (db *DCDB) AnotherDropTable(tableName string) error {
	return db.ExecSQL(`DROP TABLE "` + tableName + `"`)
}

func (db *DCDB) CreateSomeTable(tableName string, columns string) error {
	return db.ExecSQL(`CREATE SEQUENCE "` + tableName + `_id_seq" START WITH 1;
				CREATE TABLE "` + tableName + `" (
				"id" bigint NOT NULL  default nextval('` + tableName + `_id_seq'),
				` + columns + `
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER SEQUENCE "` + tableName + `_id_seq" owned by "` + tableName + `".id;
				ALTER TABLE ONLY "` + tableName + `" ADD CONSTRAINT "` + tableName + `_pkey" PRIMARY KEY (id);`)
}

func (db *DCDB) CreateNewTable(columns [][]string, tableName string, global int64, stateID string) error {
	colsSQL := ""
	colsSQL2 := ""
	sqlIndex := ""
	for _, data := range columns {
		colType := ``
		colDef := ``
		switch data[1] {
		case "text":
			colType = `varchar(102400)`
		case "int64":
			colType = `bigint`
			colDef = `NOT NULL DEFAULT '0'`
		case "time":
			colType = `timestamp`
		case "hash":
			colType = `bytea`
		case "double":
			colType = `double precision`
		case "money":
			colType = `decimal (30, 0)`
			colDef = `NOT NULL DEFAULT '0'`
		}
		colsSQL += `"` + data[0] + `" ` + colType + " " + colDef + " ,\n"
		colsSQL2 += `"` + data[0] + `": "ContractConditions(\"MainCondition\")",`
		if data[2] == "1" {
			sqlIndex += `CREATE INDEX "` + tableName + `_` + data[0] + `_index" ON "` + tableName + `" (` + data[0] + `);`
		}
	}
	colsSQL2 = colsSQL2[:len(colsSQL2)-1]

	sql := `CREATE SEQUENCE "` + tableName + `_id_seq" START WITH 1;
				CREATE TABLE "` + tableName + `" (
				"id" bigint NOT NULL  default nextval('` + tableName + `_id_seq'),
				` + colsSQL + `
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER SEQUENCE "` + tableName + `_id_seq" owned by "` + tableName + `".id;
				ALTER TABLE ONLY "` + tableName + `" ADD CONSTRAINT "` + tableName + `_pkey" PRIMARY KEY (id);`
	fmt.Println(sql)
	err := db.ExecSQL(sql)
	if err != nil {
		return err
	}

	err = db.ExecSQL(sqlIndex)
	if err != nil {
		return err
	}

	prefix := `global`
	if global == 0 {
		//if p.TxMaps.Int64["global"] == 0 {
		//prefix = p.TxStateIDStr
		prefix = stateID
	}
	err = db.ExecSQL(`INSERT INTO "`+prefix+`_tables" ( name, columns_and_permissions ) VALUES ( ?, ? )`,
		tableName, `{"general_update":"ContractConditions(\"MainCondition\")", "update": {`+colsSQL2+`},
		"insert": "ContractConditions(\"MainCondition\")", "new_column":"ContractConditions(\"MainCondition\")"}`)
	if err != nil {
		return err
	}
	return nil
}
