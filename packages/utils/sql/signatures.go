package sql

import "fmt"

func (db *DCDB) GetValueFromSignatures(tablePrefix string, name string) (string, error) {
	return db.Single(fmt.Sprintf(`select value from "%s_signatures" where name=?`, tablePrefix), name).String()
}

func (db *DCDB) GetValueConditionsFromSignatures(tablePrefix string, name string) (map[string]string, error) {
	return db.OneRow(`SELECT value, conditions FROM "`+tablePrefix+`_signatures" where name=?`, name).String()
}
