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

func (db *DCDB) IsLanguageExists(tablePrefix string, name string) (string, error) {
	return db.Single(`select name from "`+tablePrefix+"_languages"+`" where name=?`, name).String()
}

func (db *DCDB) IsLanguageExistsFromBytes(tablePrefix string, name []byte) (string, error) {
	return db.Single(`select name from "`+tablePrefix+"_languages"+`" where name=?`, name).String()
}

func (db *DCDB) CreateLanguagesTable(id string) error {
	return db.ExecSQL(`CREATE TABLE "` + id + `_languages" (
				"name" varchar(100)  NOT NULL DEFAULT '',
				"res" jsonb,
				"conditions" text  NOT NULL DEFAULT '',
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER TABLE ONLY "` + id + `_languages" ADD CONSTRAINT "` + id + `_languages_pkey" PRIMARY KEY (name);
				`)
}

func (db *DCDB) CreateFirstLanguagesRecord(id string, sid string) error {
	return db.ExecSQL(`INSERT INTO "`+id+`_languages" (name, res, conditions) VALUES
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?),
		(?, ?, ?)`,
		`dateformat`, `{"en": "YYYY-MM-DD", "ru": "DD.MM.YYYY"}`, sid,
		`timeformat`, `{"en": "YYYY-MM-DD HH:MI:SS", "ru": "DD.MM.YYYY HH:MI:SS"}`, sid,
		`Gender`, `{"en": "Gender", "ru": "Пол"}`, sid,
		`male`, `{"en": "Male", "ru": "Мужской"}`, sid,
		`female`, `{"en": "Female", "ru": "Женский"}`, sid)
}
