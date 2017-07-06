package sql

func (db *DCDB) CreateStopDaemons(time int64) error {
	return db.ExecSQL(`INSERT INTO stop_daemons(stop_time) VALUES (?)`, time)
}

func (db *DCDB) DeleteStopDaemon() error {
	return db.ExecSQL(`DELETE FROM stop_daemons`)
}

func (db *DCDB) SelectStopTime() (int64, error) {
	return db.Single(`SELECT stop_time FROM stop_daemons`).Int64()
}
