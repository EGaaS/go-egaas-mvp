package sql

func (db *DCDB) GetDataBlockIdFromRollback(rollbackID int64) (map[string]string, error) {
	return db.OneRow(`SELECT data, block_id FROM "rollback" WHERE rb_id = ?`, rollbackID).String()
}
