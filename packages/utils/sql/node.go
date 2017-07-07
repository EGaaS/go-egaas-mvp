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

func (db *DCDB) GettimeFromUpdFullNodes() (int64, error) {
	return db.Single("SELECT time FROM upd_full_nodes").Int64()
}

func (db *DCDB) GetNodeHost(nodeID string) (string, error) {
	return db.Single("SELECT host FROM full_nodes WHERE id = ?", nodeID).String()
}

func (db *DCDB) GetNodeID(stateID int64, walletID int64, delegateStateID int64, delegateWalletID int64) (int64, error) {
	return db.Single("SELECT id FROM full_nodes WHERE final_delegate_state_id = ? OR final_delegate_wallet_id = ? OR state_id = ? OR wallet_id = ?", delegateStateID, delegateWalletID, stateID, walletID).Int64()
}

func (db *DCDB) GetZeroBlockKeyID(publicKey []byte) (int64, error) {
	return db.Single(`SELECT id FROM my_node_keys WHERE block_id = 0 AND public_key = [hex]`, publicKey).Int64()
}

func (db *DCDB) CreateFullNode(walletID int64, host string) error {
	return db.ExecSQL(`INSERT INTO full_nodes (wallet_id, host) VALUES (?,?)`, walletID, host)
}
