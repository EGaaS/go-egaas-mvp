package sql

func (db *DCDB) IsSystemRestoreAccessActive(stateID int64) (map[string]int64, error) {
	return db.OneRow("SELECT active FROM system_restore_access WHERE state_id  =  ?", stateID).Int64()
}
