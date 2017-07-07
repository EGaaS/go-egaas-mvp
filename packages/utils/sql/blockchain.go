package sql

import (
	"database/sql"
	"fmt"
)

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
func (db *DCDB) GetAllDataFromBlockchain(blockID int64) (*sql.Rows, error) {
	return db.Query(db.FormatQuery(`SELECT data FROM block_chain	WHERE id > ? ORDER BY id ASC`), blockID)
}

func (db *DCDB) GetIDDataFromBlockchainLimited(startBlockID int64, limit int) (*sql.Rows, error) {
	return db.Query(db.FormatQuery("SELECT id, data FROM block_chain WHERE id > ? ORDER BY id DESC LIMIT "+fmt.Sprintf(`%d`, limit)+` OFFSET 0`), startBlockID)
}

func (db *DCDB) GetIDDataFromBlockchain(startID int64, finishID int64) (*sql.Rows, error) {
	return db.Query(db.FormatQuery(`SELECT id, data FROM block_chain	WHERE id > ? AND id <= ? ORDER BY id`),
		startID, finishID)
}

func (db *DCDB) GetLastBlockchainRecord() (map[string]string, error) {
	return db.OneRow("SELECT * FROM block_chain ORDER BY id DESC").String()
}

func (db *DCDB) DeleteBlockchainFrom(blockID int64) (int64, error) {
	return db.ExecSQLGetAffect("DELETE FROM block_chain WHERE id > ?", blockID)
}

func (db *DCDB) DeleteFromBlockchain(blockID []byte) error {
	return db.ExecSQL("DELETE FROM block_chain WHERE id = ?", blockID)
}

func (db *DCDB) DeleteFromBlockchainInt64(blockID int64) error {
	return db.ExecSQL("DELETE FROM block_chain WHERE id = ?", blockID)
}

func (db *DCDB) IsBlockchainBlockExists(blockID int64) (int64, error) {
	return db.Single("SELECT id FROM block_chain WHERE id = ?", blockID).Int64()
}

func (db *DCDB) CreateBlockchain(blockID int64, hash []byte, stateID int64, walletID int64, time int64, data []byte) (int64, error) {
	return db.ExecSQLGetAffect("INSERT INTO block_chain (id, hash, state_id, wallet_id, time, data) VALUES (?, [hex], ?, ?, ?, [hex])",
		blockID, hash, stateID, walletID, time, data)
}

func (db *DCDB) CreateBlockchainWithTx(blockID int64, hash []byte, data []byte, stateID int64, walletID int64, time int64, tx int) error {
	return db.ExecSQL("INSERT INTO block_chain (id, hash, data, state_id, wallet_id, time, tx) VALUES (?, [hex], [hex], ?, ?, ?, ?)",
		blockID, hash, data, stateID, walletID, time, tx)
}

func (db *DCDB) GetHashDataFromBlockchain(blockID int64, hash *[]byte, data *[]byte) error {
	return db.QueryRow(db.FormatQuery("SELECT hash, data FROM block_chain WHERE id  =  ?"), blockID).Scan(hash, data)
}
