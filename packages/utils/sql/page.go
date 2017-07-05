package sql

func (db *DCDB) GetPageMenus(pagePrefix, pageName string) (string, error) {
	return db.Single(`SELECT menu FROM "`+pagePrefix+`_pages" WHERE name = ?`, pageName).String()
}

func (db *DCDB) GetPage(tablePrefix, pageName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_pages" WHERE name = ?`, pageName).String()
}

func (db *DCDB) GetInterfacePages(tableprefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tableprefix+`_pages" where menu!='0' order by name`, -1)
}

func (db *DCDB) GetInterfaceBlocks(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_pages" where menu='0' order by name`, -1)
}
