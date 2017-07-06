package sql

func (db *DCDB) GetMaxBlockID() (int64, error) {
	return db.Single("select max(id) from block_chain").Int64()
}

func (db *DCDB) Get30BlocksFrom(blockID int64) ([]map[string]string, error) {
	return db.GetAll(`SELECT  b.hash, b.state_id, b.wallet_id, b.time, b.tx, b.id FROM block_chain as b
		where b.id > $1	order by b.id desc limit 30 offset 0`, -1, blockID)
}

func (db *DCDB) GetLast30Blocks() ([]map[string]string, error) {
	return db.GetAll(`SELECT  b.hash, b.state_id, b.wallet_id, b.time, b.tx, b.id FROM block_chain as b
		order by b.id desc limit 30 offset 0`, -1)
}

func (db *DCDB) GetLast30BlockInfoFrom(blockID int64) ([]map[string]string, error) {
	return db.GetAll(`SELECT b.data, b.time, b.tx, b.id FROM block_chain as b
		where b.id > $1	order by b.id desc limit 30 offset 0`, -1, blockID)
}

func (db *DCDB) GetLast30TxAndWalletInfo(blockID int64) ([]map[string]string, error) {
	return db.GetAll(`SELECT   b.wallet_id, b.time, b.tx, b.id FROM block_chain as b
			where b.id > $1	order by b.id desc limit 30 offset 0`, -1, blockID)
}

func (db *DCDB) GetHashFromBlockhain(blockID int64) (string, error) {
	return db.Single("SELECT hash FROM block_chain WHERE id = ?", blockID).String()
}

func (db *DCDB) GetDataFromBlockchain(blockID int64) ([]byte, error) {
	return db.Single("SELECT data FROM block_chain WHERE id  =  ?", blockID).Bytes()
}
