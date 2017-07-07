package sql

import "database/sql"

func (db *DCDB) GetDataBlockIdFromRollback(rollbackID int64) (map[string]string, error) {
	return db.OneRow(`SELECT data, block_id FROM "rollback" WHERE rb_id = ?`, rollbackID).String()
}

func (db *DCDB) Get1000Rollback() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM rollback ORDER BY rb_id DESC LIMIT 100`, -1)
}

func (db *DCDB) GetRollback(rollbackID int64) (map[string]string, error) {
	return db.OneRow(`select * from rollback where rb_id=?`, rollbackID).String()
}

func (db *DCDB) GetRollbackInfo(rollbackID int64) (map[string]string, error) {
	return db.OneRow(`select r.*, b.time from rollback as r
			left join block_chain as b on b.id=r.block_id
			where r.rb_id=?`, rollbackID).String()
}

func (db *DCDB) GetRollbackTx(transactionHash string) (*sql.Rows, error) {
	return db.QueryRows("SELECT tx_hash, table_name, table_id FROM rollback_tx WHERE tx_hash = [hex] ORDER BY id DESC", transactionHash)
}

func (db *DCDB) DeleteFromRollbackTx(txHash string) error {
	return db.ExecSQL("DELETE FROM rollback_tx WHERE tx_hash = [hex]", txHash)
}

func (db *DCDB) AnotherDeleteFromRollbackTx(txHash string) error {
	return db.ExecSQL(`DELETE FROM rollback_tx WHERE tx_hash = [hex] AND table_name = ?`, txHash, "system_states")
}

func (db *DCDB) InsertIntoRollback(data string, blockID int64) (string, error) {
	return db.ExecSQLGetLastInsertID("INSERT INTO rollback ( data, block_id ) VALUES ( ?, ? )", "rollback", data, blockID)
}

func (db *DCDB) CreateRollbackTX(blockID int64, hash string, tableName string, tableID string) error {
	return db.ExecSQL("INSERT INTO rollback_tx ( block_id, tx_hash, table_name, table_id ) VALUES (?, [hex], ?, ?)",
		blockID, hash, tableName, tableID)
}

func (db *DCDB) CreateRollback(data string, blockID int64) (string, error) {
	return db.ExecSQLGetLastInsertID("INSERT INTO rollback ( data, block_id ) VALUES ( ?, ? )", "rollback", data, blockID)
}

func (db *DCDB) DeleteFromRollback(rollbackID int64) error {
	return db.ExecSQL("DELETE FROM rollback WHERE rb_id = ?", rollbackID)
}

func (db *DCDB) SelectTableIDFromRollbackTx(hash string) (int64, error) {
	return db.Single(`SELECT table_id FROM rollback_tx WHERE tx_hash = [hex] AND table_name = ?`, hash, "system_states").Int64()
}
