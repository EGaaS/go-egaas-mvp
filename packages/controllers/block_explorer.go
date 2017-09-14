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
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/EGaaS/go-egaas-mvp/packages/consts"
	"github.com/EGaaS/go-egaas-mvp/packages/exchangeapi"
	"github.com/EGaaS/go-egaas-mvp/packages/lib"
	//"github.com/EGaaS/go-egaas-mvp/packages/smart"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
)

const NBlockExplorer = `block_explorer`

type TXExplorer struct {
	Hash             string
	Sender           string
	SenderAddress    string
	Recipient        string
	RecipientAddress string
	Comment          string
	Amount           string
	EGS              string
	Commission       string
	CommissionEGS    string
	ComWallet        string
	ComAddress       string
}

type blockExplorerPage struct {
	Data       *CommonPage
	List       []map[string]string
	Latest     int64
	BlockId    int64
	Public     bool
	BlockData  map[string]string
	TXs        []TXExplorer
	SinglePage int64
	History    exchangeapi.History
	TX         exchangeapi.TXInfo
}

func init() {
	newPage(NBlockExplorer)
}

func (c *Controller) BlockExplorer() (string, error) {
	pageData := blockExplorerPage{Data: c.Data}

	blockId := utils.StrToInt64(c.r.FormValue("blockId"))
	walletId := c.r.FormValue("wallet")
	hash := c.r.FormValue("hash")
	pageData.SinglePage = utils.StrToInt64(c.r.FormValue("singlePage"))
	pageData.Public = strings.HasPrefix(c.r.URL.String(), `/blockexplorer`)
	if blockId > 0 {
		pageData.BlockId = blockId
		blockInfo, err := c.OneRow(`SELECT b.* FROM block_chain as b
		where b.id=?`, blockId).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if len(blockInfo) > 0 {
			var transfer bool
			blockInfo[`hash`] = hex.EncodeToString([]byte(blockInfo[`hash`]))
			blockInfo[`size`] = utils.IntToStr(len(blockInfo[`data`]))
			if len(blockInfo[`wallet_id`]) > 0 {
				blockInfo[`wallet_address`] = lib.AddressToString(uint64(utils.StrToInt64(blockInfo[`wallet_id`])))
			} else {
				blockInfo[`wallet_address`] = ``
			}
			tmp := hex.EncodeToString([]byte(blockInfo[`data`]))
			out := ``
			for i, ch := range tmp {
				out += string(ch)
				if (i & 1) != 0 {
					out += ` `
				}
			}
			if blockId > 1 {
				parent, err := c.Single("SELECT hash FROM block_chain where id=?", blockId-1).String()
				if err == nil {
					blockInfo[`parent`] = hex.EncodeToString([]byte(parent))
				} else {
					blockInfo[`parent`] = err.Error()
				}
			}
			txlist := make([]string, 0)
			block := ([]byte(blockInfo[`data`]))[1:]
			utils.ParseBlockHeader(&block)
			//			fmt.Printf("Block OK %v sign=%d %d %x", *pblock, len((*pblock).Sign), len(block), block)
			for len(block) > 0 {
				size := int(utils.DecodeLength(&block))
				if size == 0 || len(block) < size {
					break
				}
				var name string
				itype := int(block[0])
				if itype < 128 {
					if stype, ok := consts.TxTypes[itype]; ok {
						name = stype
						if name == `DLTTransfer` {
							transfer = true
						}
					} else {
						name = fmt.Sprintf("unknown %d", itype)
					}
				} else {
					itype -= 128
					tmp := make([]byte, 4)
					for i := 0; i < itype; i++ {
						tmp[4-itype+i] = block[i+1]
					}
					//idc := int32(binary.BigEndian.Uint32(tmp))
					/*contract := smart.GetContractById(idc)
					if contract != nil {
						name = contract.Name
					} else {
						name = fmt.Sprintf(`Unknown=%d`, idc)
					}*/
				}
				txlist = append(txlist, name)
				block = block[size:]
			}
			blockInfo[`data`] = out
			blockInfo[`tx_list`] = strings.Join(txlist, `, `)
			if transfer {
				tx, _ := utils.DB.GetAll(`select table_id, tx_hash from rollback_tx where block_id=? and table_name='dlt_transactions' order by id`,
					-1, blockId)
				var txexp TXExplorer
				for _, itx := range tx {
					if txexp.Hash != hex.EncodeToString([]byte(itx[`tx_hash`])) {
						if len(txexp.Hash) > 0 {
							pageData.TXs = append(pageData.TXs, txexp)
						}
						txexp = TXExplorer{
							Hash: hex.EncodeToString([]byte(itx[`tx_hash`])),
						}
					}
					item, err := utils.DB.OneRow(`select * from dlt_transactions where id=?`, itx[`table_id`]).String()
					if err == nil && len(item) > 0 {
						if len(txexp.Sender) == 0 {
							txexp.Sender = item[`sender_wallet_id`]
							txexp.SenderAddress = lib.AddressToString(uint64(utils.StrToInt64(txexp.Sender)))
							txexp.Recipient = item[`recipient_wallet_id`]
							txexp.RecipientAddress = lib.AddressToString(uint64(utils.StrToInt64(txexp.Recipient)))
							txexp.Amount = item[`amount`]
							txexp.EGS = lib.EGSMoney(item[`amount`])
							txexp.Comment = item[`comment`]
						} else {
							txexp.Commission = item[`amount`]
							txexp.CommissionEGS = lib.EGSMoney(item[`amount`])
							txexp.ComWallet = item[`recipient_wallet_id`]
							txexp.ComAddress = lib.AddressToString(uint64(utils.StrToInt64(txexp.ComWallet)))
						}
					}
				}
				if len(txexp.Hash) > 0 {
					pageData.TXs = append(pageData.TXs, txexp)
				}
			}
		}
		pageData.BlockData = blockInfo
	} else if len(walletId) != 0 {
		pageData.SinglePage = 0
		pageData.History = exchangeapi.GetHistory(c.r)
	} else if len(hash) != 0 {
		pageData.SinglePage = 0
		pageData.TX = exchangeapi.TXStatus(c.r)
	} else {
		latest := utils.StrToInt64(c.r.FormValue("latest"))
		if latest > 0 {
			curid, _ := c.Single("select max(id) from block_chain").Int64()
			if curid <= latest {
				return ``, nil
			}
		}
		limit := `30`
		if pageData.Public {
			limit = `100`
		}
		blockExplorer, err := c.GetAll(`SELECT  b.hash, b.state_id, b.wallet_id, b.time, b.tx, b.id FROM block_chain as b
		order by b.id desc limit `+limit+` offset 0`, -1)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		for ind := range blockExplorer {
			blockExplorer[ind][`hash`] = hex.EncodeToString([]byte(blockExplorer[ind][`hash`]))
			if len(blockExplorer[ind][`wallet_id`]) > 0 {
				blockExplorer[ind][`wallet_address`] = lib.AddressToString(uint64(utils.StrToInt64(blockExplorer[ind][`wallet_id`])))
			} else {
				blockExplorer[ind][`wallet_address`] = ``
			}
			/*			if blockExplorer[ind][`tx`] == `[]` {
							blockExplorer[ind][`tx_count`] = `0`
						} else {
							var tx []string
							json.Unmarshal([]byte(blockExplorer[ind][`tx`]), &tx)
							if tx != nil && len(tx) > 0 {
								blockExplorer[ind][`tx_count`] = utils.IntToStr(len(tx))
							}
						}*/
		}
		pageData.List = blockExplorer
		if blockExplorer != nil && len(blockExplorer) > 0 {
			pageData.Latest = utils.StrToInt64(blockExplorer[0][`id`])
		}
	}
	return proceedTemplate(c, NBlockExplorer, &pageData)
}
