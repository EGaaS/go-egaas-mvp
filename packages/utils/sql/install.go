package sql

func (db *DCDB) GetInstallationState() (string, error) {
	return db.Single("SELECT progress FROM install").String()
}

func (db *DCDB) MarkInstallationComplete() error {
	return db.ExecSQL(`INSERT INTO install (progress) VALUES ('complete')`)
}
