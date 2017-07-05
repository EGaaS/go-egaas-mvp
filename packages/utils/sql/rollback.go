package sql

func (db *DCDB) GetDataBlockIdFromRollback(rollbackID int64) (map[string]string, error) {
	return db.OneRow(`SELECT data, block_id FROM "rollback" WHERE rb_id = ?`, rollbackID).String()
}

func (db *DCDB) Get1000Rollback() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM rollback ORDER BY rb_id DESC LIMIT 100`, -1)
}

func (db *DCDB) GetRollback(rollbackID int64) (map[string]string, error) {
	return db.OneRow(`select * from rollback where rb_id=?`, rollbackID).String()
}
