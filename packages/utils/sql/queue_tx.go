package sql

func (db *DCDB) CreateQueueTx(hash []byte, data string) error {
	return db.ExecSQL("INSERT INTO queue_tx (hash, data) VALUES ([hex], [hex])", hash, data)
}

func (db *DCDB) CreateAnotherQueueTx(hash []byte, data []byte) error {
	return db.ExecSQL("INSERT INTO queue_tx (hash, data) VALUES ([hex], [hex])", hash, data)
}

func (db *DCDB) IsTransactionQueueExists(txHash []byte) (int64, error) {
	return db.Single("SELECT count(hash) FROM queue_tx WHERE hex(hash) = ?", txHash).Int64()
}

func (db *DCDB) CreateGateQueueTx(hash []byte, hex []byte) error {
	return db.ExecSQL(`INSERT INTO queue_tx (hash, data, from_gate) VALUES ([hex], [hex], 1)`, hash, hex)
}

func (db *DCDB) DeleteFromQueueTx(hash []byte) error {
	return db.ExecSQL(`DELETE FROM queue_tx WHERE hex(hash) = ?`, hash)
}
