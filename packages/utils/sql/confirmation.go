package sql

func (db *DCDB) IsBlockConfirmationExists(blockID int64) (int64, error) {
	return db.Single("SELECT block_id FROM confirmations WHERE block_id= ?", blockID).Int64()
}

func (db *DCDB) MarkConfirmations(good int64, bad int64, time int64, blockID int64) error {
	return db.ExecSQL("UPDATE confirmations SET good = ?, bad = ?, time = ? WHERE block_id = ?", good, bad, time, blockID)
}

func (db *DCDB) CreateConfirmation(good int64, bad int64, time int64, blockID int64) error {
	return db.ExecSQL("INSERT INTO confirmations ( block_id, good, bad, time ) VALUES ( ?, ?, ?, ? )", blockID, good, bad, time)
}
