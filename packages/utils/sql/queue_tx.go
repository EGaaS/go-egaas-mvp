package sql

func (db *DCDB) CreateQueueTx(hash []byte, data string) error {
	return db.ExecSQL("INSERT INTO queue_tx (hash, data) VALUES ([hex], [hex])", hash, data)
}

func (db *DCDB) CreateAnotherQueueTx(hash []byte, data []byte) error {
	return db.ExecSQL("INSERT INTO queue_tx (hash, data) VALUES ([hex], [hex])", hash, data)
}
