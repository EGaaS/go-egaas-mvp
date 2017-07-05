package sql

func (db *DCDB) GetMainLock() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM main_lock`, -1)
}
