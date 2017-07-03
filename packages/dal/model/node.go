package model

import "github.com/EGaaS/go-egaas-mvp/packages/dal/types"

type Node struct {
	Model
	//block_generator line 99
	MyStateID  types.DALInt64
	MyWalletID types.DALInt64
	//block_generator line 111
	DelegateWalletID types.DALInt64
	DelegateStateID  types.DALInt64
}

func (t *Node) SetDatasource(ds map[string]types.DataSource) error {
	for field, source := range ds {
		switch field {
		case "MyStateID":
			t.MyStateID.SetDatasource(source)
		case "MyWalletID":
			t.MyWalletID.SetDatasource(source)
		case "DelegateWalletID":
			t.DelegateWalletID.SetDatasource(source)
		case "DelegateStateID":
			t.DelegateStateID.SetDatasource(source)
		default:
			return UnknownDatasorceField
		}
	}
	return nil
}

var NodeDataSource = map[string]types.DataSource{
	"MyStateID":        types.DataSource{ResName: "config|system_recognized_states", ParamName: "state_id"},
	"MyWalletID":       types.DataSource{ResName: "config", ParamName: "wallet_id"},
	"DelegateWalletID": types.DataSource{ResName: "system_recognized_states", ParamName: "delegate_wallet_id"},
	"DelegateStateID":  types.DataSource{ResName: "system_recognized_states", ParamName: "delegate_state_id"},
}
