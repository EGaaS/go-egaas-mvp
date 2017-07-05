package sql

func (db *DCDB) GetValueFromMenu(pagePrefix, menuName string) (string, error) {
	return db.Single(`SELECT value FROM "`+pagePrefix+`_menu" WHERE name = ?`, menuName).String()
}

func (db *DCDB) GetMenu(pagePrefix, menuName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+pagePrefix+`_menu" WHERE name = ?`, menuName).String()
}

func (db *DCDB) GetAllMenus(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_menu" order by name`, -1)
}
