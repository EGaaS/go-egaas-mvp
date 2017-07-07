package sql

func (db *DCDB) GetOrderedSmartContracts(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_smart_contracts" order by id`, -1)
}

func (db *DCDB) GetSmartContractsByID(tablePrefix string, contractID int64) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_smart_contracts" WHERE id = ?`, contractID).String()
}

func (db *DCDB) GetSmartContractsByName(tablePrefix string, contractName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_smart_contracts" WHERE name = ?`, contractName).String()
}

func (db *DCDB) GetSmartContractID(tablePrefix string, contractName string) (string, error) {
	return db.Single(`SELECT id FROM "`+tablePrefix+`_smart_contracts" WHERE name = ?`, contractName).String()
}

func (db *DCDB) GetSmartContractIDFromByte(tablePrefix string, contractName []byte) (int64, error) {
	return db.Single(`SELECT id FROM "`+tablePrefix+`_smart_contracts" WHERE name = ?`, contractName).Int64()
}

func (db *DCDB) IsSmartContractActive(tablePrefix string, contractID string) (string, error) {
	return db.Single(`SELECT active FROM "`+tablePrefix+`_smart_contracts" WHERE id = ?`, contractID).String()
}

func (db *DCDB) IsSmartContractActiveAndID(tablePrefix string, contractID string) (map[string]string, error) {
	return db.OneRow(`SELECT id, active FROM "`+tablePrefix+`_smart_contracts" WHERE id = ?`, contractID).String()
}

func (db *DCDB) GetConditionFromSmartContract(tablePrefix string, contractID string) (string, error) {
	return db.Single(`SELECT conditions FROM "`+tablePrefix+`_smart_contracts" WHERE id = ?`, contractID).String()
}

func (db *DCDB) GetIdConditionsActiveFromSmartContract(tablePrefix string, contractName string) (map[string]string, error) {
	return db.OneRow(`SELECT id,conditions, active FROM "`+tablePrefix+`_smart_contracts" WHERE name = ?`, contractName).String()
}

func (db *DCDB) CreateSmartContractTable(id string) error {
	return db.ExecSQL(`CREATE SEQUENCE "` + id + `_smart_contracts_id_seq" START WITH 1;
				CREATE TABLE "` + id + `_smart_contracts" (
				"id" bigint NOT NULL  default nextval('` + id + `_smart_contracts_id_seq'),
				"name" varchar(100)  NOT NULL DEFAULT '',
				"value" text  NOT NULL DEFAULT '',
				"wallet_id" bigint  NOT NULL DEFAULT '0',
				"active" character(1) NOT NULL DEFAULT '0',
				"conditions" text  NOT NULL DEFAULT '',
				"variables" bytea  NOT NULL DEFAULT '',
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER SEQUENCE "` + id + `_smart_contracts_id_seq" owned by "` + id + `_smart_contracts".id;
				ALTER TABLE ONLY "` + id + `_smart_contracts" ADD CONSTRAINT "` + id + `_smart_contracts_pkey" PRIMARY KEY (id);
				CREATE INDEX "` + id + `_smart_contracts_index_name" ON "` + id + `_smart_contracts" (name);
				`)
}

func (db *DCDB) CreateSmartContractMainCondition(id string, walletID int64) error {
	return db.ExecSQL(`INSERT INTO "`+id+`_smart_contracts" (name, value, wallet_id, active) VALUES
		(?, ?, ?, ?)`,
		`MainCondition`, `contract MainCondition {
            data {}
            conditions {
                    if(StateVal("gov_account")!=$citizen)
                    {
                        warning "Sorry, you don't have access to this action."
                    }
            }
            action {}
    }`, walletID, 1,
	)
}

func (db *DCDB) UpdateSmartContractConditions(id string, conditions string) error {
	return db.ExecSQL(`UPDATE "`+id+`_smart_contracts" SET conditions = ?`, conditions)
}
