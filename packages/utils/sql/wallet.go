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
