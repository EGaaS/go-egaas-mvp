package sql

import "database/sql"

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

func (db *DCDB) IsVerifiedTransactionExists(hash []byte) (int64, error) {
	return db.Single("SELECT count(hash) FROM transactions WHERE hex(hash) = ? and verified = 1", hash).Int64()
}

func (db *DCDB) DeleteBadTransactions() (int64, error) {
	return db.ExecSQLGetAffect("DELETE FROM transactions WHERE verified = 0 AND used = 0 AND counter > 10")
}

func (db *DCDB) GetHashDataOfUnsendedTransactions() (*sql.Rows, error) {
	return db.Query("SELECT hash, data FROM transactions WHERE sent  =  0")
}

func (db *DCDB) MarkTransactionSended(transactionHash []byte) (int64, error) {
	return db.ExecSQLGetAffect("UPDATE transactions SET sent = 1 WHERE hex(hash) = ?", transactionHash)
}

func (db *DCDB) MarkTransactionUsed(transactionHash []byte) (int64, error) {
	return db.ExecSQLGetAffect("UPDATE transactions SET used=1 WHERE hex(hash) = ?", transactionHash)
}

func (db *DCDB) GetTransactionData(transactionHash []byte) ([]byte, error) {
	return db.Single("SELECT data FROM transactions WHERE hex(hash) = ?", transactionHash).Bytes()
}

func (db *DCDB) MarkTransactionUnverifiedExec() error {
	return db.ExecSQL("UPDATE transactions SET verified = 0 WHERE verified = 1 AND used = 0")
}

func (db *DCDB) MarkTransactionUnverified() (int64, error) {
	return db.ExecSQLGetAffect("UPDATE transactions SET verified = 0 WHERE verified = 1 AND used = 0")
}

func (db *DCDB) MarkTransactionUnusedAndUnverified(transactionHash string) (int64, error) {
	return db.ExecSQLGetAffect("UPDATE transactions SET used=0, verified = 0 WHERE hex(hash) = ?", transactionHash)
}

func (db *DCDB) DeleteUsedTransactions() (int64, error) {
	return db.ExecSQLGetAffect("DELETE FROM transactions WHERE used = 1")
}

func (db *DCDB) MarkTransactionStatusByBlockIDBytes(blockID int64, transactionHash []byte) error {
	return db.ExecSQL("UPDATE transactions_status SET block_id = ? WHERE hex(hash) = ?", blockID, transactionHash)
}

func (db *DCDB) MarkTransactionStatusByBlockID(blockID int64, transactionHash string) error {
	return db.ExecSQL("UPDATE transactions_status SET block_id = ? WHERE hex(hash) = ?", blockID, transactionHash)
}

func (db *DCDB) MarkTransactionStatusError(blockHash string, errorText string) error {
	return db.ExecSQL("UPDATE transactions_status SET error = ? WHERE hex(hash) = ?", errorText, blockHash)
}

func (db *DCDB) DeleteFromLogTransaction(transactionHash string) (int64, error) {
	return db.ExecSQLGetAffect("DELETE FROM log_transactions WHERE hex(hash) = ?", transactionHash)
}

func (db *DCDB) DeleteTransaction(hash string) (int64, error) {
	return db.ExecSQLGetAffect("DELETE FROM transactions WHERE hex(hash) = ?", hash)
}

func (db *DCDB) DeleteTransactionByte(hash []byte) (int64, error) {
	return db.ExecSQLGetAffect("DELETE FROM transactions WHERE hex(hash) = ?", hash)
}

func (db *DCDB) SetTransactionStatusError(errorMessage string, hash []byte) error {
	return db.ExecSQL("UPDATE transactions_status SET error = ? WHERE hex(hash) = ?", errorMessage, hash)
}

func (db *DCDB) GetTransactionCounter(hash []byte) (int64, error) {
	return db.Single("SELECT counter FROM transactions WHERE hex(hash) = ?", hash).Int64()
}

func (db *DCDB) CreateTransaction(hash []byte, data []byte, forSelfUse int, transactionType int64, walletID int64, citizenID int64, thirdVar int64, counter int64) error {
	return db.ExecSQL(`INSERT INTO transactions (hash, data, for_self_use, type, wallet_id, citizen_id, third_var, counter, verified) VALUES ([hex], [hex], ?, ?, ?, ?, ?, ?, 1)`,
		hash, data, forSelfUse, transactionType, walletID, citizenID, thirdVar, counter)
}

func (db *DCDB) CreateDltTransaction(senderWalletID int64, recepientWalletID int64, recepientWalletAddress string,
	amount string, comission string, comment []byte, time int64, blockID int64) (string, error) {
	return db.ExecSQLGetLastInsertID(`INSERT INTO dlt_transactions ( sender_wallet_id, recipient_wallet_id, recipient_wallet_address, amount, commission, comment, time, block_id ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )`, "dlt_transactions",
		senderWalletID, recepientWalletID, recepientWalletAddress, amount, comission, comment, time, blockID)
}

func (db *DCDB) CreateAnotherDltTransaction(senderWalletID int64, recepientWalletID int64, recepientWalletAddress string,
	amount string, comission int, comment string, time int64, blockID int64) (string, error) {
	return db.ExecSQLGetLastInsertID(`INSERT INTO dlt_transactions ( sender_wallet_id, recipient_wallet_id, recipient_wallet_address, amount, commission, comment, time, block_id ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )`, "dlt_transactions",
		senderWalletID, recepientWalletID, recepientWalletAddress, amount, comission, comment, time, blockID)
}

func (db *DCDB) DeleteUnusedAndUnverifiedByHash(hash []byte) (int64, error) {
	return db.ExecSQLGetAffect("DELETE FROM transactions WHERE hex(hash) = ? AND verified=0 AND used = 0", hash)
}

func (db *DCDB) GetAllDataHashFromTransactionsAndQueue() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM (SELECT data, hash FROM queue_tx UNION SELECT data, hash FROM transactions WHERE verified = 0 AND used = 0)  AS x`, -1)
}

func (db *DCDB) CreateIncorrectTransactionTx(time int64, hash []byte, errorText string) error {
	return db.ExecSQL(`INSERT INTO incorrect_tx (time, hash, err) VALUES (?, [hex], ?)`, time, hash, errorText)
}

func (db *DCDB) GetHashFromLogTransactions(hash []byte) (string, error) {
	return db.Single(`SELECT hash FROM log_transactions WHERE hex(hash) = ?`, hash).String()
}
