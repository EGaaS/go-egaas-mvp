package sql

func (db *DCDB) GetPageMenus(pagePrefix, pageName string) (string, error) {
	return db.Single(`SELECT menu FROM "`+pagePrefix+`_pages" WHERE name = ?`, pageName).String()
}

func (db *DCDB) GetPage(tablePrefix, pageName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_pages" WHERE name = ?`, pageName).String()
}
