package sql

func (db *DCDB) GetTxCountBySenderOrRecepient(senderWalletID, recipientWalletID int64, recipientWalletAddress string) (int64, error) {
	return db.Single(`SELECT count(id) FROM dlt_transactions where sender_wallet_id = ? OR
		                       recipient_wallet_id = ? OR recipient_wallet_address = ?`, senderWalletID, recipientWalletID, recipientWalletAddress).Int64()
}

func (db *DCDB) GetAllTxBySenderOrRecepient(senderWalletID, recipientWalletID int64, recipientWalletAddress string, limit string) ([]map[string]string, error) {
	return db.GetAll(`SELECT d.*, w.wallet_id as sw, wr.wallet_id as rw FROM dlt_transactions as d
		        left join dlt_wallets as w on w.wallet_id=d.sender_wallet_id
		        left join dlt_wallets as wr on wr.wallet_id=d.recipient_wallet_id
				where sender_wallet_id=? OR 
		        recipient_wallet_id=?  OR
		        recipient_wallet_address=? order by d.id desc  `+limit, -1, senderWalletID, senderWalletID, recipientWalletAddress)
}

func (db *DCDB) CreateTxStatus(hash []byte, time int64, txType int32, walletID int64, citizenID int64) error {
	return db.ExecSQL(`INSERT INTO transactions_status (
			hash, time,	type, wallet_id, citizen_id	) VALUES (
			[hex], ?, ?, ?, ? )`, hash, time, txType, walletID, citizenID)
}

func (db *DCDB) GetBlockIDError(hash string) (map[string]string, error) {
	return db.OneRow(`SELECT block_id, error FROM transactions_status WHERE hash = [hex]`, hash).String()
}

func (db *DCDB) GetTransactionByRecipient(walletID int64) (map[string]string, error) {
	return db.OneRow(`select * from dlt_transactions where recipient_wallet_id=? order by id`, walletID).String()
}

func (db *DCDB) GetLast100Transactions() ([]map[string]string, error) {
	return db.GetAll("SELECT hex(hash) as hex_hash, verified, used, high_rate, for_self_use, user_id, third_var, counter, sent FROM transactions", 100)
}

func (db *DCDB) IsTransactionLogExists(txHash []byte) (int64, error) {
	return db.Single("SELECT count(hash) FROM log_transactions WHERE hex(hash) = ?", txHash).Int64()
}

func (db *DCDB) IsTransactionExists(txHash []byte) (int64, error) {
	return db.Single("SELECT count(hash) FROM transactions WHERE hex(hash) = ?", txHash).Int64()
}
