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
	//	"encoding/json"

	"github.com/EGaaS/go-egaas-mvp/packages/converter"
)

const aExplorer = `ajax_explorer`

// ExplorerJSON is a structure for the answer of ajax_explorer ajax request
type ExplorerJSON struct {
	Data   []map[string]string `json:"data"`
	Latest int64               `json:"latest"`
}

func init() {
	newPage(aExplorer, `json`)
}

// AjaxExplorer is a controller of ajax_explorer request
func (c *Controller) AjaxExplorer() interface{} {
	result := ExplorerJSON{}
	latest := converter.StrToInt64(c.r.FormValue("latest"))
	if latest > 0 {
		result.Latest, _ = c.GetMaxBlockID()
		if result.Latest > latest {
			explorer, err := c.Get30BlocksFrom(latest)
			if err == nil {
				for ind := range explorer {
					explorer[ind][`hash`] = hex.EncodeToString([]byte(explorer[ind][`hash`]))
					if len(explorer[ind][`wallet_id`]) > 0 {
						explorer[ind][`wallet_address`] = converter.AddressToString(converter.StrToInt64(explorer[ind][`wallet_id`]))
					} else {
						explorer[ind][`wallet_address`] = ``
					}
				}
				result.Data = explorer
				if explorer != nil && len(explorer) > 0 {
					result.Latest = converter.StrToInt64(explorer[0][`id`])
				}
			}
		}
	}
	if result.Data == nil {
		result.Data = make([]map[string]string, 0)
	}
	return result
}
