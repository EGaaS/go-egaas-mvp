package sql

func (db *DCDB) CreateNodeKeys(privateKey string, publicKey string) error {
	return db.ExecSQL(`INSERT INTO my_node_keys (private_key, public_key, block_id) VALUES (?, [hex], ?)`, privateKey, publicKey, 1)
}

func (db *DCDB) CreateNodeKeysWithoutBlockID(privateKey []byte, publicKey []byte) error {
	return db.ExecSQL(`INSERT INTO my_node_keys (public_key, private_key) VALUES ([hex], ?)`, publicKey, privateKey)
}

func (db *DCDB) GetAllUpdFullNodes() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM upd_full_nodes`, -1)
}

func (db *DCDB) GetAllFullNodes() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM rollback ORDER BY rb_id DESC LIMIT 100`, -1)
}
