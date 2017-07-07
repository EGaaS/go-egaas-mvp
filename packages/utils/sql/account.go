package sql

import "fmt"

func (db *DCDB) GetAccountAmount(stateID, citizenID int64) (string, error) {
	return db.Single(fmt.Sprintf(`select amount from "%d_accounts" where citizen_id=?`,
		stateID), citizenID).String()
}

// TODO anonims deprecated
func (db *DCDB) GetAnonyms(stateid, citizenID int64) ([]map[string]string, error) {
	return db.GetAll(fmt.Sprintf(`select anon.*, acc.amount from "%d_anonyms" as anon
	left join "%[1]d_accounts" as acc on acc.citizen_id=anon.id_anonym
	where anon.id_citizen=?`, stateid), -1, citizenID)
}

func (db *DCDB) CreateAnonyms(stateID int64, citizenID int64, accountID int64, encryptedKey string) error {
	return db.ExecSQL(fmt.Sprintf(`INSERT INTO "%d_anonyms" (id_citizen, id_anonym, encrypted)
			VALUES (?,?,[hex])`, stateID), citizenID, accountID, encryptedKey)
}

func (db *DCDB) CreateAnonymsTable(id string) error {
	return db.ExecSQL(`CREATE TABLE "` + id + `_anonyms" (
				"id_citizen" bigint NOT NULL DEFAULT '0',
				"id_anonym" bigint NOT NULL DEFAULT '0',
				"encrypted" bytea  NOT NULL DEFAULT ''
				);
				CREATE INDEX "` + id + `_anonyms_index_id" ON "` + id + `_anonyms" (id_citizen);`)
}
