// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"fmt"

	"github.com/EGaaS/go-egaas-mvp/packages/converter"
	"github.com/EGaaS/go-egaas-mvp/packages/template"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
	"github.com/EGaaS/go-egaas-mvp/packages/utils/sql"
)

var (
	isGlobal bool
)

/*
Adding state tables should be spelled out in state settings
*/

// NewStateInit initializes NewState transaction
func (p *Parser) NewStateInit() error {

	fields := []map[string]string{{"state_name": "string"}, {"currency_name": "string"}, {"public_key": "bytes"}, {"sign": "bytes"}}
	err := p.GetTxMaps(fields)
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

// NewStateGlobal checks if the state or the currency exists
func (p *Parser) NewStateGlobal(country, currency string) error {
	if !isGlobal {
		list, err := sql.DB.GetAllTables()
		if err != nil {
			return err
		}
		isGlobal = converter.InSliceString(`global_currencies_list`, list) && converter.InSliceString(`global_states_list`, list)
	}
	if isGlobal {
		if id, err := sql.DB.GetStateID(country); err != nil {
			return err
		} else if id > 0 {
			return fmt.Errorf(`State %s already exists`, country)
		}
		if id, err := sql.DB.GetCurrencyID(currency); err != nil {
			return err
		} else if id > 0 {
			return fmt.Errorf(`Currency %s already exists`, currency)
		}
	}
	return nil
}

// NewStateFront checks conditions of NewState transaction
func (p *Parser) NewStateFront() error {
	err := p.generalCheck(`new_state`)
	if err != nil {
		return p.ErrInfo(err)
	}

	// Check InputData
	verifyData := map[string]string{"state_name": "state_name", "currency_name": "currency_name"}
	err = p.CheckInputData(verifyData)
	if err != nil {
		return p.ErrInfo(err)
	}

	forSign := fmt.Sprintf("%s,%s,%d,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxWalletID, p.TxMap["state_name"], p.TxMap["currency_name"])
	CheckSignResult, err := utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false)
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return p.ErrInfo("incorrect sign")
	}
	country := string(p.TxMap["state_name"])
	if exist, err := p.IsState(country); err != nil {
		return p.ErrInfo(err)
	} else if exist > 0 {
		return fmt.Errorf(`State %s already exists`, country)
	}

	err = p.NewStateGlobal(country, string(p.TxMap["currency_name"]))
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

// NewStateMain creates state tables in the database
func (p *Parser) NewStateMain(country, currency string) (id string, err error) {
	id, err = p.ExecSQLGetLastInsertID(`INSERT INTO system_states DEFAULT VALUES`, "system_states")
	if err != nil {
		return
	}
	err = p.CreateRollbackTX(p.BlockData.BlockId, p.TxHash, "system_states", id)
	if err != nil {
		return
	}

	err = p.CreateStateTable(id)
	if err != nil {
		return
	}
	sid := "ContractConditions(`MainCondition`)" //`$citizen == ` + utils.Int64ToStr(p.TxWalletID) // id + `_citizens.id=` + utils.Int64ToStr(p.TxWalletID)
	psid := sid                                  //fmt.Sprintf(`Eval(StateParam(%s, "main_conditions"))`, id) //id+`_state_parameters.main_conditions`
	err = p.CreateStateConditions(id, sid, psid, currency, country, p.TxWalletID)
	if err != nil {
		return
	}
	err = p.CreateSmartContractTable(id)
	if err != nil {
		return
	}
	err = p.CreateSmartContractMainCondition(id, p.TxWalletID)

	if err != nil {
		return
	}

	err = p.UpdateSmartContractConditions(id, sid)
	if err != nil {
		return
	}

	err = p.CreateTables(id)
	if err != nil {
		return
	}

	err = p.CreateTablesRecords(id, sid, psid)
	if err != nil {
		return
	}

	err = p.CreatePagesTable(id)
	if err != nil {
		return
	}

	err = p.CreateFirstPagesRecords(id, sid)
	if err != nil {
		return
	}

	err = p.CreateMenuTable(id)
	if err != nil {
		return
	}
	err = p.CreateFirstMenuRecord(id, sid)
	if err != nil {
		return
	}

	err = p.CreateCitizensTable(id)
	if err != nil {
		return
	}

	pKey, err := p.GetSingleWalletPublicKeyBytes(p.TxWalletID)
	if err != nil {
		return
	}

	err = p.CreateFirstCitizenRecord(id, p.TxWalletID, converter.BinToHex(pKey))
	if err != nil {
		return
	}
	err = p.CreateLanguagesTable(id)
	if err != nil {
		return
	}
	err = p.CreateFirstLanguagesRecord(id, sid)
	if err != nil {
		return
	}

	err = p.CreateSignaturesTable(id)
	if err != nil {
		return
	}

	err = p.CreateAppsTable(id)
	if err != nil {
		return
	}

	err = p.CreateAnonymsTable(id)
	if err != nil {
		return
	}

	err = template.LoadContract(id)
	return
}

// NewState proceeds NewState transaction
func (p *Parser) NewState() error {
	var pkey string
	country := string(p.TxMap["state_name"])
	currency := string(p.TxMap["currency_name"])
	id, err := p.NewStateMain(country, currency)
	if err != nil {
		return p.ErrInfo(err)
	}
	if isGlobal {
		_, err = p.selectiveLoggingAndUpd([]string{"gstate_id", "state_name", "timestamp date_founded"},
			[]interface{}{id, country, p.BlockData.Time}, "global_states_list", nil, nil, true)

		if err != nil {
			return p.ErrInfo(err)
		}
		_, err = p.selectiveLoggingAndUpd([]string{"currency_code", "settings_table"},
			[]interface{}{currency, id + `_state_parameters`}, "global_currencies_list", nil, nil, true)
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	if pkey, err = p.GetSingleWalletPublicKey(p.TxWalletID); err != nil {
		return p.ErrInfo(err)
	} else if len(p.TxMaps.Bytes["public_key"]) > 30 && len(pkey) == 0 {
		_, err = p.selectiveLoggingAndUpd([]string{"public_key_0"}, []interface{}{converter.HexToBin(p.TxMaps.Bytes["public_key"])}, "dlt_wallets",
			[]string{"wallet_id"}, []string{converter.Int64ToStr(p.TxWalletID)}, true)
	}
	return err
}

// NewStateRollback rollbacks NewState transaction
func (p *Parser) NewStateRollback() error {
	id, err := p.SelectTableIDFromRollbackTx(p.TxHash)
	if err != nil {
		return p.ErrInfo(err)
	}
	err = p.autoRollback()
	if err != nil {
		return p.ErrInfo(err)
	}

	for _, name := range []string{`menu`, `pages`, `citizens`, `languages`, `signatures`, `tables`,
		`smart_contracts`, `state_parameters`, `apps`, `anonyms` /*, `citizenship_requests`*/} {
		err = p.DropTable(id, name)
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	err = p.AnotherDeleteFromRollbackTx(p.TxHash)
	if err != nil {
		return p.ErrInfo(err)
	}

	maxID, err := p.GetMaxSystemStateID()
	if err != nil {
		return p.ErrInfo(err)
	}
	// обновляем AI
	// update  the AI
	err = p.SetAI("system_states", maxID+1)
	if err != nil {
		return p.ErrInfo(err)
	}
	err = p.DeleteFromSystemStates(id)
	if err != nil {
		return p.ErrInfo(err)
	}

	return nil
}
