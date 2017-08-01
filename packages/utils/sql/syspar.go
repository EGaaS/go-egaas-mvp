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

package sql

import (
	"sync"

	"github.com/EGaaS/go-egaas-mvp/packages/converter"

	"github.com/shopspring/decimal"
)

type sysParName string

const (
	// NumberNodes is the number of nodes
	NumberNodes = sysParName(`number_of_dlt_nodes`)
	// FuelRate is the rate
	FuelRate = sysParName(`fuel_rate`)
	// OpPrice is the costs of operations
	OpPrice = sysParName(`op_price`)
	// GapsBetweenBlocks is the time between blocks
	GapsBetweenBlocks = sysParName(`gaps_between_blocks`)
	// BlockchainURL is the address of the blockchain file.  For those who don't want to collect it from nodes
	BlockchainURL = sysParName(`blockchain_url`)
	// MaxBlockSize is the maximum size of the block
	MaxBlockSize = sysParName(`max_block_size`)
	// MaxTxSize is the maximum size of the transaction
	MaxTxSize = sysParName(`max_tx_size`)
	// MaxTxCount is the maximum count of the transactions
	MaxTxCount = sysParName(`max_tx_count`)
	// MaxColumns is the maximum columns in tables
	MaxColumns = sysParName(`max_columns`)
	// MaxIndexes is the maximum indexes in tables
	MaxIndexes = sysParName(`max_indexes`)
	// MaxBlockUserTx is the maximum number of user's transactions in one block
	MaxBlockUserTx = sysParName(`max_block_user_tx`)
	// UpdFullNodesPeriod is the maximum number of user's transactions in one block
	UpdFullNodesPeriod = sysParName(`upd_full_nodes_period`)
	// RecoveryAddress is the recovery address
	RecoveryAddress = sysParName(`recovery_address`)
	// CommissionWallet is the address for commissions
	CommissionWallet = sysParName(`commission_wallet`)
)

var (
	cache = map[sysParName]string{
		BlockchainURL: "https://raw.githubusercontent.com/egaas-blockchain/egaas-blockchain.github.io/master/testnet_blockchain",
		// For compatible of develop versions
		// Remove later
		GapsBetweenBlocks:  `3`,
		MaxBlockSize:       `67108864`,
		MaxTxSize:          `33554432`,
		MaxTxCount:         `100000`,
		MaxColumns:         `50`,
		MaxIndexes:         `10`,
		MaxBlockUserTx:     `100`,
		UpdFullNodesPeriod: `3600`, // 3600 is for the test time, then we have to put 86400`
		RecoveryAddress:    `8275283526439353759`,
		CommissionWallet:   `8275283526439353759`,
	}
	mutex = &sync.Mutex{}
)

// SysUpdate reloads/updates values of system parameters
func SysUpdate() error {
	list, err := DB.GetAll(`select * from system_parameters`, -1)
	if err != nil {
		return err
	}
	mutex.Lock()
	for _, item := range list {
		for key, value := range item {
			cache[sysParName(key)] = value
		}
	}
	mutex.Unlock()
	return nil
}

// SysDecimal returns big integer value
func SysDecimal(name sysParName) (decimal.Decimal, error) {
	return decimal.NewFromString(string(name))
}

// SysInt64 returns int64 value of the system parameter
func SysInt64(name sysParName) int64 {
	return converter.StrToInt64(SysString(name))
}

// SysInt returns int64 value of the system parameter
func SysInt(name sysParName) int {
	return converter.StrToInt(SysString(name))
}

// SysString returns string value of the system parameter
func SysString(name sysParName) string {
	mutex.Lock()
	ret := cache[name]
	mutex.Unlock()
	return ret
}
