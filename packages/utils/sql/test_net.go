package sql

func (db *DCDB) GetPrivateKeyFromTestNet(walletID int64) (string, error) {
	return db.Single(`select private from testnet_emails where wallet=?`, walletID).String()
}

func (db *DCDB) CreateTestNetKeys(sitizenID, stateID int64, spriv string, idnew int64) error {
	return db.ExecSQL(`insert into testnet_keys (id, state_id, private, wallet) values(?,?,?,?)`,
		stateID, stateID, spriv, idnew)
}

func (db *DCDB) GetAllTestnetKeys(citizenID, stateID int64) (int64, error) {
	return db.Single(`select count(id) from testnet_keys where id=? and state_id=?`, citizenID, stateID).Int64()
}

func (db *DCDB) GetAvailableTestnetKeys(citizenID, stateID int64) (int64, error) {
	return db.Single(`select count(id) from testnet_keys where id=? and state_id=? and status=0`, citizenID, stateID).Int64()
}

func (db *DCDB) SelectMainDataFromTestnet(ID int64) (map[string]string, error) {
	return db.OneRow(`select country,currency,wallet, private from testnet_emails where id=?`, ID).String()
}

func (db *DCDB) UpdateTestnetEmails(walletAddress int64, privateKey string, ID int64) error {
	return db.ExecSQL(`update testnet_emails set wallet=?, private=? where id=?`, walletAddress, privateKey, ID)
}
