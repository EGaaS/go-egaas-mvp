package sql

func (db *DCDB) GetShortInfoBlock() (map[string]int64, error) {
	return db.OneRow("SELECT state_id, wallet_id, block_id, time, hex(hash) as hash FROM info_block").Int64()
}

func (db *DCDB) GetInfoBlockHash() (string, error) {
	return db.Single("SELECT hex(hash) as hash FROM info_block").String()
}
