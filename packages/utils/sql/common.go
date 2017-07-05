package sql

func (db *DCDB) GetColumnsCount(tableName string) (int64, error) {
	return db.Single("SELECT count(column_name) FROM information_schema.columns WHERE table_name=?", tableName).Int64()
}

func (db *DCDB) GetColumnType(tableName, columnName string) (map[string]string, error) {
	return db.OneRow(`select data_type,character_maximum_length from information_schema.columns 
		where table_name = ? and column_name = ?`, tableName, columnName).String()
}

func (db *DCDB) DropTables() error {
	return db.ExecSQL(`
	DO $$ DECLARE
	    r RECORD;
	BEGIN
	    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
		EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
	    END LOOP;
	END $$;
	`)
}
