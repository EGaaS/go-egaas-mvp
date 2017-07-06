package sql

import "fmt"

func (db *DCDB) GetCitizenshipPrice(stateID string) (int64, error) {
	return db.Single(`SELECT value FROM "` + stateID + `_state_parameters" where name='citizenship_price'`).Int64()
}

func (db *DCDB) GetEGSRate(stateID string) (float64, error) {
	return db.Single(`SELECT value FROM "`+stateID+`_state_parameters" WHERE name = ?`, `egs_rate`).Float64()
}

func (db *DCDB) GetStateCoords(stateID string) (string, error) {
	return db.Single(`SELECT coords FROM "` + stateID + `_state_details"`).String()
}

func (db *DCDB) GetAllSystemStatesIDs() ([]string, error) {
	return db.GetList(`SELECT id FROM system_states`).String()
}

func (db *DCDB) GetAllSystemStatesIDsOrdered() ([]map[string]string, error) {
	return db.GetAll(`select id from system_states order by id`, -1)
}

func (db *DCDB) GetStateValue(stateID string, stateName string) (string, error) {
	return db.Single(fmt.Sprintf(`SELECT value FROM "%s_state_parameters" WHERE name = ?`, stateID), stateName).String()
}

func (db *DCDB) GetStateParameterByParameter(stateID string, parameter string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+stateID+`_state_parameters" WHERE parameter = ?`, parameter).String()
}

func (db *DCDB) GetStateParameterByName(stateID, name string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+stateID+`_state_parameters" WHERE name = ?`, name).String()
}

func (db *DCDB) GetAllStateParameters(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_state_parameters" order by name`, -1)
}

func (db *DCDB) GetParametersFromEA() ([]string, error) {
	return db.GetList(`SELECT parameter FROM ea_state_parameters`).String()
}

func (db *DCDB) GetStateNames(statePrefix string) ([]string, error) {
	return db.GetList(`SELECT name FROM "` + statePrefix + `_state_parameters"`).String()
}

func (db *DCDB) GetGlobalStateID(stateName string) (int64, error) {
	return db.Single("select gstate_id from global_states_list where state_name=?", stateName).Int64()
}

func (db *DCDB) GetGlobalStateName(stateID int64) (string, error) {
	return db.Single(`select state_name from global_states_list where gstate_id=?`, stateID).String()
}

func (db *DCDB) GetStateID(stateName string) (int64, error) {
	return DB.Single(`select id from global_states_list where state_name=?`, stateName).Int64()
}

func (db *DCDB) GetEAStateLaws() ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM ea_state_laws`, -1)
}

func (db *DCDB) GetEAStateParameters() ([]string, error) {
	return db.GetList(`SELECT parameter FROM ea_state_parameters`).String()
}
