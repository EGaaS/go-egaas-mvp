package model

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestSelectQuery(t *testing.T) {
	transactions := DltTransaction()
	query := transactions.Read(transactions.Amount, transactions.Comission).Query()
	if query != "select amount, comission from dlt_transactions;" {
		t.Error("select error. query: ", query)
	}
}

func TestInsertQuery(t *testing.T) {
	transactions := DltTransaction()
	amount, _ := decimal.NewFromString("10")
	comission, _ := decimal.NewFromString("12")
	query := transactions.Create(
		transactions.Amount.SetValue(amount),
		transactions.Comission.SetValue(comission)).Query()
	if query != "insert into dlt_transactions (amount, comission) values (10, 12);" {
		t.Error("insert error. query: ", query)
	}
}

func TestWhereQuery(t *testing.T) {
	transactions := DltTransaction()

	query := transactions.Read(
		transactions.Amount, transactions.Comission).Where(
		Condition{transactions.Amount, Less, "10"}).And(
		Condition{transactions.Comission, GreaterOrEqual, "15"}).Or(
		Condition{transactions.Comission, Equal, "14"}).Query()
	if query != "select amount, comission from dlt_transactions where amount < 10 and comission >= 15 or comission = 14;" {
		t.Error("where error. query: ", query)
	}
}

func TestDeleteQuery(t *testing.T) {
	transactions := DltTransaction()

	query := transactions.Delete(Condition{transactions.Amount, Less, "10"}).And(
		Condition{transactions.Comission, GreaterOrEqual, "15"}).Query()

	if query != "delete from dlt_transactions where amount < 10 and comission >= 15;" {
		t.Error("delete error. query: ", query)
	}
}
