package sql

func (db *DCDB) GetCurrencyID(currency string) (int64, error) {
	return db.Single(`select id from global_currencies_list where currency_code=?`, currency).Int64()
}
