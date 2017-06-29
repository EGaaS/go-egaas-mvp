package model

import "github.com/EGaaS/go-egaas-mvp/packages/dal/types"

type Transaction struct {
	ID        types.DALString
	Hash      types.DALByteArray
	Data      types.DALByteArray
	Used      types.DALint8
	HignRate  types.DALint8
	Type      types.DALint8
	WalletID  types.DALInt64
	CitizenID types.DALInt64
	Counter   types.DALint8
	Sent      types.DALint8
}
