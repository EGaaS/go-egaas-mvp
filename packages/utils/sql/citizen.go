package sql

func (db *DCDB) GetOrderedSitizensIDs(stateID string, address int64) ([]map[string]string, error) {
	return db.GetAll(`select id from "`+stateID+`_citizens" where id>=? order by id`, 7, address)
}

func (db *DCDB) IsCitizenExist(stateID string, citizenID int64) (int64, error) {
	return db.Single(`select id from "`+stateID+`_citizens" where id=?`, citizenID).Int64()
}

func (db *DCDB) GetCitizenshipRequests(stateID string, walletID int64) (map[string]int64, error) {
	return db.OneRow(`select id, approved from "`+stateID+`_citizenship_requests" where dlt_wallet_id=? order by id desc`,
		walletID).Int64()
}

func (db *DCDB) GetCitizenshipRequestsFull(stateID string, walletID int64) (map[string]string, error) {
	return db.OneRow(`select * from "`+stateID+`_citizenship_requests" where dlt_wallet_id=? order by id desc`,
		walletID).String()
}

func (db *DCDB) GetCitizenName(stateID string, citizenID int64) (string, error) {
	return db.Single(`SELECT name FROM "`+stateID+`_citizens" WHERE id = ?`, citizenID).String()
}

func (db *DCDB) GetCitizenAvatar(stateID string, citizenID int64) (string, error) {
	return db.Single(`SELECT avatar FROM "`+stateID+`_citizens" WHERE id = ?`, citizenID).String()
}

func (db *DCDB) CreateCitizensTable(id string) error {
	return db.ExecSQL(`CREATE TABLE "` + id + `_citizens" (
				"id" bigint NOT NULL DEFAULT '0',
				"public_key_0" bytea  NOT NULL DEFAULT '',				
				"block_id" bigint NOT NULL DEFAULT '0',
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER TABLE ONLY "` + id + `_citizens" ADD CONSTRAINT "` + id + `_citizens_pkey" PRIMARY KEY (id);
				`)
}

func (db *DCDB) CreateFirstCitizenRecord(id string, walletID int64, privateKey []byte) error {
	return db.ExecSQL(`INSERT INTO "`+id+`_citizens" (id,public_key_0) VALUES (?, [hex])`, walletID, privateKey)
}
