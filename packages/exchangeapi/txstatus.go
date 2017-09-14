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
	Confirmations    string `json:"confirmations"`
	Hash             string `json:"txhash"`
	Amount           string `json:"amount"`
	EGS              string `json:"egs"`
	Time             string `json:"time"`
	Sender           string `json:"sender"`
	Recipient        string `json:"recipient"`
	SenderAddress    string `json:"sender_address"`
	RecipientAddress string `json:"recipient_address"`
	Comment          string
	Commission       string
	CommissionEGS    string
	ComWallet        string
	ComAddress       string
	Error            string `json:"error"`
}

func TXStatus(r *http.Request) TXInfo {
	return txstatus(r).(TXInfo)
}

func txstatus(r *http.Request) interface{} {
	var (
		result TXInfo
	)

	/*	tx, err := utils.DB.OneRow(`SELECT block_id, error, time FROM transactions_status WHERE hash = [hex]`, r.FormValue(`hash`)).String()
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
		result.BlockId = tx["block_id"]
		result.Time = tx["time"]*/
	result.Hash = r.FormValue(`hash`)
	tx, err := utils.DB.GetAll(`select block_id, table_id from rollback_tx where tx_hash=[hex] and table_name='dlt_transactions'
							order by id`, -1, r.FormValue(`hash`))
	if err != nil {
		result.Error = err.Error()
		return result
	}
	if len(tx) > 0 {
		txitem, err := utils.DB.OneRow(`select * from dlt_transactions	where id=?`,
			tx[0]["table_id"]).String()
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.BlockId = txitem[`block_id`]
		result.Sender = txitem[`sender_wallet_id`]
		result.SenderAddress = lib.AddressToString(uint64(utils.StrToInt64(result.Sender)))
		result.Recipient = txitem[`recipient_wallet_id`]
		result.RecipientAddress = lib.AddressToString(uint64(utils.StrToInt64(result.Recipient)))
		result.Amount = txitem[`amount`]
		result.EGS = lib.EGSMoney(txitem[`amount`])
		result.Comment = txitem[`comment`]
		result.Time = txitem[`time`]
	}
	if len(tx) > 1 {
		txitem, err := utils.DB.OneRow(`select * from dlt_transactions	where id=?`,
			tx[1]["table_id"]).String()
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.Commission = txitem[`amount`]
		result.CommissionEGS = lib.EGSMoney(txitem[`amount`])
		result.ComWallet = txitem[`recipient_wallet_id`]
		result.ComAddress = lib.AddressToString(uint64(utils.StrToInt64(result.ComWallet)))
	}

	/*	txitem, err := utils.DB.OneRow(`select * from dlt_transactions
			 where id=?`,
			tx["block_id"], list[0][`table_id`], list[1][`table_id`]).String()
		if err != nil {
			result.Error = err.Error()
			return result
		}*/
	conf, err := utils.DB.OneRow(`select * from confirmations where block_id=?`,
		result.BlockId).String()
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
	return result
}
