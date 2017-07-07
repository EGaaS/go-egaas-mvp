package sql

func (db *DCDB) CreateConfig(firstLoad, firstLoadBlockchainURL string, autoReload int32) error {
	return db.ExecSQL("INSERT INTO config (first_load_blockchain, first_load_blockchain_url, auto_reload) VALUES (?, ?, ?)", firstLoad, firstLoadBlockchainURL, autoReload)
}

func (db *DCDB) UpdateConfig(dltWalletID int64) error {
	return db.ExecSQL(`UPDATE config SET dlt_wallet_id = ?`, dltWalletID)
}

func (db *DCDB) UpdateConfigBlockID(blockID int64) error {
	return db.ExecSQL(`UPDATE config SET my_block_id = ?`, blockID)
}

func (db *DCDB) AnotherUpdateBlockID(newBlockID int64, whereBlockID int64) error {
	return db.ExecSQL("UPDATE config SET my_block_id = ? WHERE my_block_id < ?", newBlockID, whereBlockID)
}

func (db *DCDB) GetFirstLoadBlockchainURL() (string, error) {
	return db.Single("SELECT first_load_blockchain_url FROM config").String()
}

func (db *DCDB) GetFirstLoadBlockchain() (string, error) {
	return db.Single("SELECT first_load_blockchain FROM config").String()
}

func (db *DCDB) SetCurrentLoadBlockchainFile() error {
	return db.ExecSQL(`UPDATE config SET current_load_blockchain = 'file'`)
}

func (db *DCDB) MarkLoadBlockainNodes() error {
	return db.ExecSQL(`UPDATE config SET current_load_blockchain = 'nodes'`)
}

func (db *DCDB) GetBadBlocksFromConfig() ([]byte, error) {
	return db.Single("SELECT bad_blocks FROM config").Bytes()
}

func (db *DCDB) SetBlockIDInConfig(blockID int64) error {
	return db.ExecSQL(`UPDATE config SET my_block_id = ?`, blockID)
}
