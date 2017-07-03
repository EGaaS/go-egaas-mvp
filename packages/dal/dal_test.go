package dal

import (
	"testing"

	"github.com/EGaaS/go-egaas-mvp/packages/dal/model"
)

func TestSelectQuery(t *testing.T) {
	transaction := &model.Transaction{}
	transaction.SetDatasource(model.TransactionDatasource)
	transaction.Hash.Set([]byte("123123313"))
	transaction.CitizenID.Set(123)
	transaction.Counter.Set(1)
	transaction.WalletID.Set(2)
}
