package sql

func (db *DCDB) GetLanguageRes(tablePrefix string, langName string) (string, error) {
	return db.Single(`SELECT res FROM "`+tablePrefix+`_languages" where name=?`, langName).String()
}
