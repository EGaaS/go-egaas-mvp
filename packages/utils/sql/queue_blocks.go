package sql

func (db *DCDB) CreateQueueBlock(newDataHash []byte, fullNodeID int64, newDataBlockID int64) error {
	return db.ExecSQL(`INSERT INTO queue_blocks (hash, full_node_id, block_id) VALUES ([hex], ?, ?) ON CONFLICT DO NOTHING`,
		newDataHash, fullNodeID, newDataBlockID)
}

func (db *DCDB) GetOneBlockQueue() (map[string]string, error) {
	return db.OneRow("SELECT * FROM queue_blocks").String()
}
