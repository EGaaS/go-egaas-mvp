package sql

import "fmt"

func (db *DCDB) GetLanguageRes(tablePrefix string, langName string) (string, error) {
	return db.Single(`SELECT res FROM "`+tablePrefix+`_languages" where name=?`, langName).String()
}

func (db *DCDB) GetLanguages(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT name, res FROM "`+tablePrefix+`_languages" order by name`, -1)
}

func (db *DCDB) GetAllLanguages(stateID int) ([]map[string]string, error) {
	return db.GetAll(fmt.Sprintf(`select * from "%d_languages"`, stateID), -1)
}
