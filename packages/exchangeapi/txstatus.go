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

package exchangeapi

import (
	"net/http"

	"github.com/EGaaS/go-egaas-mvp/packages/lib"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
)

type TXInfo struct {
	BlockId          string `json:"block_id"`
	Confirmations    string  `json:"confirmations"`
	Hash             string `json:"txhash"`
	Amount           string `json:"amount"`
	EGS              string `json:"egs"`
	Time             string `json:"time"`
	Sender           string `json:"sender"`
	Recipient        string `json:"recipient"`
	SenderAddress    string `json:"sender_address"`
	RecipientAddress string `json:"recipient_address"`
	Error            string `json:"error"`
}

func txstatus(r *http.Request) interface{} {
	var (
		result TXInfo
	)

	tx, err := utils.DB.OneRow(`SELECT block_id, error, time FROM transactions_status WHERE hash = [hex]`, r.FormValue(`hash`)).String()
	if err != nil {
		result.Error = err.Error()
		return result
	}
	if len(tx["error"]) > 0 {
		result.Error = tx["error"]
		return result
	}
	if tx["block_id"] == "0" || tx["block_id"] == "" {
		result.BlockId = "0"
		return result
	}
	result.Hash = r.FormValue(`hash`)
	result.BlockId = tx["block_id"]
	result.Time = tx["time"]
	list, err := utils.DB.GetAll(`select table_id from rollback_tx where block_id=? and table_name='dlt_wallets' and
						    tx_hash=[hex] order by id`,
		-1, tx["block_id"], r.FormValue(`hash`))
	if len(list) < 2 {
		return result
	}
	txitem, err := utils.DB.OneRow(`select sender_wallet_id, recipient_wallet_id, amount from dlt_transactions
		 where block_id=? and sender_wallet_id =? and recipient_wallet_id=?`,
		tx["block_id"], list[0][`table_id`], list[1][`table_id`]).String()
	if err != nil {
		result.Error = err.Error()
		return result
	}
	conf, err := utils.DB.OneRow(`select * from confirmations where block_id=?`,
		tx["block_id"]).String()
	if err != nil {
		result.Error = err.Error()
		return result
	}
	bad := float64(utils.StrToInt64(conf[`bad`]))
	good := float64(utils.StrToInt64(conf[`good`]))
	if good/bad > (good+bad)*0.51 {
		max, err := utils.DB.Single(`select max(id) from block_chain`).Int64()
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.Confirmations = utils.Int64ToStr(max - utils.StrToInt64(result.BlockId))
	}

	if len(txitem[`amount`]) > 0 {
		result.Amount = txitem[`amount`]
		result.EGS = lib.EGSMoney(result.Amount)
		result.Sender = txitem[`sender_wallet_id`]
		result.Recipient = txitem[`recipient_wallet_id`]
		result.SenderAddress = lib.AddressToString(uint64(utils.StrToInt64(result.Sender)))
		result.RecipientAddress = lib.AddressToString(uint64(utils.StrToInt64(result.Recipient)))
	}
	return result
}
