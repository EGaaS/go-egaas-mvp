package sql

func (db *DCDB) GetMigrationHistoryVersion() (string, error) {
	return db.Single(`SELECT version FROM migration_history ORDER BY id DESC LIMIT 1`).String()
}

func (db *DCDB) CreateMigrationHistory(version string, time int64) error {
	return db.ExecSQL(`INSERT INTO migration_history (version, date_applied) VALUES (?, ?)`, version, time)
}
