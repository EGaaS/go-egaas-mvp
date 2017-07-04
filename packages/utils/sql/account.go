package sql

import "fmt"

func (db *DCDB) GetAccountAmount(stateID, citizenID int64) (string, error) {
	return db.Single(fmt.Sprintf(`select amount from "%d_accounts" where citizen_id=?`,
		stateID), citizenID).String()
}

func (db *DCDB) GetAnonyms(stateid, citizenID int64) ([]map[string]string, error) {
	return db.GetAll(fmt.Sprintf(`select anon.*, acc.amount from "%d_anonyms" as anon
	left join "%[1]d_accounts" as acc on acc.citizen_id=anon.id_anonym
	where anon.id_citizen=?`, stateid), -1, citizenID)
}

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
