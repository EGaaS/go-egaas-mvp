package sql

func (db *DCDB) GetPrivateKeyFromTestNet(walletID int64) (string, error) {
	return db.Single(`select private from testnet_emails where wallet=?`, walletID).String()
}

func (db *DCDB) GetPrivateWalletFromTestNet(walletID string) (map[string]string, error) {
	return db.OneRow(`select private, wallet from testnet_emails where id=?`, walletID).String()
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

func (db *DCDB) GetStatePendingKeysCount(stateID int64) (int64, error) {
	return db.Single(`select count(id) from testnet_keys where state_id=? and status=0`, stateID).Int64()
}

func (db *DCDB) GetIdWalletPrivateKeyFromTestNet(stateID int64) (map[string]string, error) {
	return db.OneRow(`select id, wallet, private from testnet_keys where state_id=? and status=0`, stateID).String()
}

func (db *DCDB) SelectMainDataFromTestnet(ID int64) (map[string]string, error) {
	return db.OneRow(`select country,currency,wallet, private from testnet_emails where id=?`, ID).String()
}

func (db *DCDB) UpdateTestnetEmails(walletAddress int64, privateKey string, ID int64) error {
	return db.ExecSQL(`update testnet_emails set wallet=?, private=? where id=?`, walletAddress, privateKey, ID)
}

func (db *DCDB) MarkTestnetKeysStatus1(keyID int64, stateID int64, walletID int64) error {
	return db.ExecSQL(`update testnet_keys set status=1 where id=? and state_id=? and status=0 and wallet=?`,
		keyID, stateID, walletID)
}

func (db *DCDB) GetTestNetEmailID(email string, country string, currency string) (int64, error) {
	return db.Single(`select id from testnet_emails where email=? and country = ? and currency=?`,
		email, country, currency).Int64()
}

func (db *DCDB) CreateTestNetEmails(email string, country string, currency string) (string, error) {
	return db.ExecSQLGetLastInsertID(`insert into testnet_emails (email,country,currency) 
				values(?,?,?)`, `testnet_emails`, email, country, currency)
}

func (db *DCDB) CreateTestNetEmailsTable() error {
	return db.ExecSQL(`CREATE SEQUENCE testnet_emails_id_seq START WITH 1;
CREATE TABLE "testnet_emails" (
"id" integer NOT NULL DEFAULT nextval('testnet_emails_id_seq'),
"email" varchar(128) NOT NULL DEFAULT '',
"country" varchar(128) NOT NULL DEFAULT '',
"currency" varchar(32) NOT NULL DEFAULT '',
"private" varchar(64) NOT NULL DEFAULT '',
"wallet" bigint NOT NULL DEFAULT '0',
"status" integer NOT NULL DEFAULT '0',
"code" integer NOT NULL DEFAULT '0',
"validate" integer NOT NULL DEFAULT '0'
);
ALTER SEQUENCE testnet_emails_id_seq owned by testnet_emails.id;
ALTER TABLE ONLY "testnet_emails" ADD CONSTRAINT testnet_emails_pkey PRIMARY KEY (id);
CREATE INDEX testnet_index_email ON "testnet_emails" (email);`)
}

func (db *DCDB) CreateTestNetKeysTable() error {
	return db.ExecSQL(`CREATE TABLE "testnet_keys" (
		"id" bigint NOT NULL DEFAULT '0',
		"state_id" integer NOT NULL DEFAULT '0',
		"private" varchar(64) NOT NULL DEFAULT '',
		"wallet" bigint NOT NULL DEFAULT '0',
		"status" integer NOT NULL DEFAULT '0'
		);
		CREATE INDEX testnet_index_keys ON "testnet_keys" (id,state_id,status);`)
}
