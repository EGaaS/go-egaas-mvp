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

package controllers

import (
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
)

const NSystemInfo = `system_info`

type systemInfoPage struct {
	Data      *CommonPage
	List      []map[string]string
	Latest    int64
	BlockId   int64
	UpdFullNodes []map[string]string
	MainLock []map[string]string
	Rollback []map[string]string
	FullNodes []map[string]string
	Votes []map[string]string
	Confirmations []map[string]string
	Wallets []map[string]string
	WalletsTransactions []map[string]string
}

func init() {
	newPage(NSystemInfo)
}

func (c *Controller) SystemInfo() (string, error) {
	var err error
	pageData := systemInfoPage{Data: c.Data}

	pageData.UpdFullNodes, err = c.GetAll(`SELECT * FROM upd_full_nodes`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	pageData.MainLock, err = c.GetAll(`SELECT * FROM main_lock`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	pageData.Rollback, err = c.GetAll(`SELECT * FROM rollback ORDER BY rb_id DESC LIMIT 100`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	pageData.FullNodes, err = c.GetAll(`SELECT * FROM full_nodes`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	pageData.Votes, err = c.GetAll(`SELECT address_vote, sum(amount) as sum FROM dlt_wallets WHERE address_vote !='' GROUP BY address_vote ORDER BY sum(amount) DESC LIMIT 10`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	pageData.Confirmations, err = c.GetAll(`SELECT * FROM confirmations ORDER BY block_id DESC LIMIT 100`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	pageData.Wallets, err = c.GetAll(`SELECT wallet_id, amount, address_vote,	host,	last_forging_data_upd, hex(node_public_key),	hex(public_key_0),	hex(public_key_1),	hex(public_key_2),	rb_id FROM dlt_wallets ORDER BY amount DESC LIMIT 100`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	pageData.WalletsTransactions, err = c.GetAll(`SELECT id, amount,	block_id,	commission,rb_id,	recipient_wallet_address,	recipient_wallet_id,	sender_wallet_id, time  FROM dlt_transactions ORDER BY block_id DESC LIMIT 100`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	return proceedTemplate(c, NSystemInfo, &pageData)
}
