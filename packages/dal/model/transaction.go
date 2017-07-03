package model

import "github.com/EGaaS/go-egaas-mvp/packages/dal/types"

type Transaction struct {
	Model
	Hash      types.DALByteArray
	Data      types.DALByteArray
	Used      types.DALInt8
	HighRate  types.DALInt8
	Type      types.DALInt8
	WalletID  types.DALInt64
	CitizenID types.DALInt64
	Counter   types.DALInt8
	Sent      types.DALInt8
}

func NewTransaction() *Transaction {
	result := &Transaction{}
	result.CitizenID.SetNull(true)
	result.Counter.SetNull(true)
	result.Data.SetNull(true)
	result.Hash.SetNull(true)
	result.HighRate.SetNull(true)
	result.ID.SetNull(true)
	result.LastInsertID.SetNull(true)
	result.Sent.SetNull(true)
	result.Type.SetNull(true)
	result.Used.SetNull(true)
	result.WalletID.SetNull(true)
	return result
}

func (t *Transaction) SetDatasource(ds map[string]types.DataSource) error {
	for field, source := range ds {
		switch field {
		case "hash":
			t.Hash.SetDatasource(source)
		case "data":
			t.Data.SetDatasource(source)
		case "used":
			t.Used.SetDatasource(source)
		case "high_rate":
			t.HighRate.SetDatasource(source)
		case "type":
			t.Type.SetDatasource(source)
		case "wallet_id":
			t.WalletID.SetDatasource(source)
		case "citizen_id":
			t.CitizenID.SetDatasource(source)
		case "counter":
			t.Counter.SetDatasource(source)
		case "sent":
			t.Sent.SetDatasource(source)
		default:
			return UnknownDatasorceField
		}
	}
	return nil
}

var TransactionDatasource = map[string]types.DataSource{
	"hash":       types.DataSource{ResName: "transactions", ParamName: "hash"},
	"data":       types.DataSource{ResName: "transactions", ParamName: "data"},
	"used":       types.DataSource{ResName: "transactions", ParamName: "used"},
	"high_rate":  types.DataSource{ResName: "transactions", ParamName: "high_rate"},
	"type":       types.DataSource{ResName: "transactions", ParamName: "type"},
	"wallet_id":  types.DataSource{ResName: "transactions", ParamName: "wallet_id"},
	"citizen_id": types.DataSource{ResName: "transactions", ParamName: "citizen_id"},
	"counter":    types.DataSource{ResName: "transactions", ParamName: "counter"},
	"sent":       types.DataSource{ResName: "transactions", ParamName: "sent"},
}
