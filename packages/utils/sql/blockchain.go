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
