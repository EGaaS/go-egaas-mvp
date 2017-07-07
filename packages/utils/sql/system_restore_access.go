package sql

func (db *DCDB) IsSystemRestoreAccessActive(stateID int64) (map[string]int64, error) {
	return db.OneRow("SELECT active FROM system_restore_access WHERE state_id  =  ?", stateID).Int64()
}

func (db *DCDB) AnotherIsSystemRestoreAccessActive(stateID int64) (int64, error) {
	return db.Single("SELECT active FROM system_restore_access WHERE state_id = ?", stateID).Int64()
}

func (db *DCDB) GetCloseFromsystemRestoreAccess(userID int64, stateID int64) (int64, error) {
	return db.Single("SELECT close FROM system_restore_access WHERE user_id  =  ? AND state_id = ?", userID, stateID).Int64()
}

func (db *DCDB) GetAllFromSystemRestoreAccess(stateID uint32) (map[string]int64, error) {
	return db.OneRow("SELECT * FROM system_restore_access WHERE state_id  =  ?", stateID).Int64()
}

func (db *DCDB) AnotherGetAllFromSystemRestoreAccess(stateID int64) (map[string]int64, error) {
	return db.OneRow("SELECT * FROM system_restore_access WHERE state_id  =  ?", stateID).Int64()
}

func (db *DCDB) GetCitizenIDFromSystemRestoreAccess(stateID int64) (string, error) {
	return db.Single(`SELECT citizen_id FROM system_restore_access WHERE state_id = ?`, stateID).String()
}
