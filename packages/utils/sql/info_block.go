package sql

func (db *DCDB) GetShortInfoBlock() (map[string]int64, error) {
	return db.OneRow("SELECT state_id, wallet_id, block_id, time, hex(hash) as hash FROM info_block").Int64()
}

func (db *DCDB) GetInfoBlockHash() (string, error) {
	return db.Single("SELECT hex(hash) as hash FROM info_block").String()
}

func (db *DCDB) GetBlockIDFromInfoBlock() (int64, error) {
	return db.Single("SELECT block_id FROM info_block").Int64()
}

func (db *DCDB) GetUnsendedBlockIDHashFromInfoBlock() (map[string][]byte, error) {
	return db.OneRow("SELECT block_id, hash FROM info_block WHERE sent  =  0").Bytes()
}

func (db *DCDB) MarkInfoBlockSended() error {
	return db.ExecSQL("UPDATE info_block SET sent = 1")
}

func (db *DCDB) GetOneInfoBlock() (map[string]string, error) {
	return db.OneRow("SELECT * FROM info_block").String()
}

func (db *DCDB) UpdateInfoBlock(hash []byte, blockID int64, time int64, walletID int64, stateID int64) error {
	return db.ExecSQL("UPDATE info_block SET hash = [hex], block_id = ?, time = ?, wallet_id = ?, state_id = ?, sent = 0", hash, blockID, time, walletID, stateID)
}

func (db *DCDB) AnotherUpdateInfoBlock(hash []byte, blockID int64, time int64, walletID int64, stateID int64) error {
	return db.ExecSQL("UPDATE info_block SET hash = [hex], block_id = ?, time = ?, wallet_id = ?, state_id = ?",
		hash, blockID, time, walletID, stateID)
}

func (db *DCDB) CreateInfoBlock(hash []byte, blockID int64, time int64, stateID int64, walletID int64, currentVersion string) error {
	return db.ExecSQL("INSERT INTO info_block (hash, block_id, time, state_id, wallet_id, current_version) VALUES ([hex], ?, ?, ?, ?, ?)",
		hash, blockID, time, stateID, walletID, currentVersion)
}

func (db *DCDB) MySQLiteGetInfoBlockData(blockHash *[]byte, blockID *int64, blockTime *int64) error {
	return db.QueryRow("SELECT LOWER(HEX(hash)) as hash, block_id, time FROM info_block").Scan(&blockHash, &blockID, &blockTime)
}

func (db *DCDB) PgGetInfoBlockData(blockHash *[]byte, blockID *int64, blockTime *int64) error {
	return db.QueryRow("SELECT encode(hash, 'HEX')  as hash, block_id, time FROM info_block").Scan(&blockHash, &blockID, &blockTime)
}
