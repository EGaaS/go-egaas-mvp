package sql

import "fmt"

func (db *DCDB) GetValueFromSignatures(tablePrefix string, name string) (string, error) {
	return db.Single(fmt.Sprintf(`select value from "%s_signatures" where name=?`, tablePrefix), name).String()
}

func (db *DCDB) GetValueConditionsFromSignatures(tablePrefix string, name string) (map[string]string, error) {
	return db.OneRow(`SELECT value, conditions FROM "`+tablePrefix+`_signatures" where name=?`, name).String()
}

func (db *DCDB) GetConditionsFromSignatures(tablePrefix string, conditionName string) (string, error) {
	return db.Single(`SELECT conditions FROM "`+tablePrefix+`_signatures" WHERE name = ?`, conditionName).String()
}

func (db *DCDB) GetSignature(tablePrefix string, name string) (map[string]string, error) {
	return db.OneRow(fmt.Sprintf(`select * from "%s_signatures" where name=?`, tablePrefix), name).String()
}

func (db *DCDB) GetAllSignatures(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_signatures" order by name`, -1)
}

func (db *DCDB) IsSignatureExists(tableName string, name []byte) (string, error) {
	return db.Single(`select name from "`+tableName+"_signatures"+`" where name=?`, name).String()
}

func (db *DCDB) CreateSignaturesTable(id string) error {
	return db.ExecSQL(`CREATE TABLE "` + id + `_signatures" (
				"name" varchar(100)  NOT NULL DEFAULT '',
				"value" jsonb,
				"conditions" text  NOT NULL DEFAULT '',
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER TABLE ONLY "` + id + `_signatures" ADD CONSTRAINT "` + id + `_signatures_pkey" PRIMARY KEY (name);
				`)
}
