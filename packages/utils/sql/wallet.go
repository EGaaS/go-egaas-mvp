package sql

func (db *DCDB) GetWallet(walletID int64) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "dlt_wallets" WHERE wallet_id = ?`, walletID).String()
}

func (db *DCDB) GetWalletAmount(walletID int64) (int64, error) {
	return db.Single("select amount from dlt_wallets where wallet_id=?", walletID).Int64()
}

func (db *DCDB) GetWalletAmountString(walletID int64) (string, error) {
	return db.Single("select amount from dlt_wallets where wallet_id=?", walletID).String()
}

func (db *DCDB) GetWalletsIDs(address int64) ([]map[string]string, error) {
	return db.GetAll(`select wallet_id as id from dlt_wallets where wallet_id>=? order by wallet_id`, 7, address)
}

func (db *DCDB) GetWalletPublickKey(walletID string) (map[string]string, error) {
	return db.OneRow("SELECT public_key_0 FROM dlt_wallets WHERE wallet_id = ?", walletID).String()
}

func (db *DCDB) GetSingleWalletPublicKey(walletID int64) (string, error) {
	return db.Single(`select public_key_0 from dlt_wallets where wallet_id=?`, walletID).String()
}

func (db *DCDB) GetSingleWalletPublicKeyBytes(walletID int64) ([]byte, error) {
	return db.Single(`select public_key_0 from dlt_wallets where wallet_id=?`, walletID).Bytes()
}

func (db *DCDB) IsWalletExist(walletID string) (int64, error) {
	return db.Single(`select wallet_id from dlt_wallets where wallet_id=?`, walletID).Int64()
}

func (db *DCDB) GetWalletHostData(walletID int64) (map[string]string, error) {
	return db.OneRow("SELECT host, address_vote, fuel_rate FROM dlt_wallets WHERE wallet_id = ?", walletID).String()
}

func (db *DCDB) GetShortHostData(walletID int64) (map[string]string, error) {
	return db.OneRow("SELECT host, address_vote as addressVote  FROM dlt_wallets WHERE wallet_id = ?", walletID).String()
}

func (db *DCDB) GetVotes() ([]map[string]string, error) {
	return db.GetAll(`SELECT address_vote, sum(amount) as sum FROM dlt_wallets WHERE address_vote !='' GROUP BY address_vote ORDER BY sum(amount) DESC LIMIT 10`, -1)
}

func (db *DCDB) GetWalletAmountAndRollbackID(walletID int64) (map[string]string, error) {
	return db.OneRow(`select amount, rb_id from dlt_wallets where wallet_id=?`, walletID).String()
}
