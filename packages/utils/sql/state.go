package sql

import "fmt"

func (db *DCDB) GetCitizenshipPrice(stateID string) (int64, error) {
	return db.Single(`SELECT value FROM "` + stateID + `_state_parameters" where name='citizenship_price'`).Int64()
}

func (db *DCDB) GetStateCoords(stateID string) (string, error) {
	return db.Single(`SELECT coords FROM "` + stateID + `_state_details"`).String()
}

func (db *DCDB) GetSystemStates() ([]string, error) {
	return db.GetList(`SELECT id FROM system_states`).String()
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

func (db *DCDB) GetParametersFromEA() ([]string, error) {
	return db.GetList(`SELECT parameter FROM ea_state_parameters`).String()
}

func (db *DCDB) GetStateNames(statePrefix string) ([]string, error) {
	return db.GetList(`SELECT name FROM "` + statePrefix + `_state_parameters"`).String()
}
