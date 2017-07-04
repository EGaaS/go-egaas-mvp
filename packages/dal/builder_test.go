package dal

import (
	"testing"

	"github.com/EGaaS/go-egaas-mvp/packages/dal/builder"
	"github.com/EGaaS/go-egaas-mvp/packages/dal/model"
)

func createTestNode() *model.Node {
	node := &model.Node{}
	node.SetDatasource(model.NodeDataSource)
	node.MyStateID.Set(123)
	node.MyWalletID.Set(456)
	node.DelegateStateID.Set(567)
	node.DelegateWalletID.Set(321)
	return node
}
func TestPgBuilderCreate(t *testing.T) {
	node := createTestNode()

	build := builder.PgBuilder{}
	request := build.Create(&node.MyStateID, &node.MyWalletID, &node.DelegateStateID, &node.DelegateWalletID).Compile()
	if len(request.Querys) != 2 {
		t.Error("incorrect querys count. real length: ", len(request.Querys))
	}

	if request.Querys[0] != "insert into config (state_id, wallet_id) values (123, 456);" &&
		request.Querys[1] != "insert into config (state_id, wallet_id) values (123, 456);" {
		t.Error("incorrect first query: ", request.Querys[0])
	}

	if request.Querys[0] != "insert into system_recognized_states (state_id, delegate_state_id, delegate_wallet_id) values (123, 567, 321);" &&
		request.Querys[1] != "insert into system_recognized_states (state_id, delegate_state_id, delegate_wallet_id) values (123, 567, 321);" {
		t.Error("incorrect second query: ", request.Querys[1])
	}
}

func TestPgBuilderRead(t *testing.T) {
	node := createTestNode()
	build := builder.PgBuilder{}

	reader := build.Read(&node.MyStateID, &node.MyWalletID, &node.DelegateStateID, &node.DelegateWalletID).
		Where(builder.Condition{
			Field:      &node.MyStateID,
			Comparator: builder.Equal,
			Value:      node.MyStateID.String()}).
		And(builder.Condition{
			Field:      &node.DelegateStateID,
			Comparator: builder.NotEqual,
			Value:      "12"}).
		Compile()

		myStateID, myWalletID, err := d.GetMyStateIDAndWalletID()
		logger.Debug("%v", myWalletID)
		if err != nil {
			d.dbUnlock()
			logger.Error("%v", err)
			if d.dSleep(d.sleepTime) {
				break BEGIN
			}
			continue
		}

		if myStateID > 0 {
			delegate, err := d.CheckDelegateCB(myStateID)
			if err != nil {
				d.dbUnlock()
				logger.Error("%v", err)
				if d.dSleep(d.sleepTime) {
					break BEGIN
				}
				continue
			}

	if len(build.Querys) != 2 {
		t.Error("incorrect querys count. real length: ", len(reader.Querys))
	}

	if reader.Querys[0] != "select state_id, wallet_id from config;" &&
		reader.Querys[1] != "select state_id, wallet_id from config;" {
		t.Error("incorrect first query: ", reader.Querys[0])
	}

	if reader.Querys[0] != "select delegate_state_id, delegate_wallet_id from system_recognized_states where state_id = 123 and delegate_state_id != 12;" &&
		reader.Querys[1] != "select delegate_state_id, delegate_wallet_id from system_recognized_states where state_id = 123 and delegate_state_id != 12;" {
		t.Error("incorrect second query: ", build.Querys[1])
	}
}

func TestPgBuilderDelete(t *testing.T) {
	node := createTestNode()
	build := builder.PgBuilder{}

	deleter := build.Delete(&node.MyStateID, &node.MyWalletID, &node.DelegateStateID, &node.DelegateWalletID).
		Where(builder.Condition{
			Field:      &node.MyStateID,
			Comparator: builder.Equal,
			Value:      node.MyStateID.String()}).
		Or(builder.Condition{
			Field:      &node.DelegateStateID,
			Comparator: builder.NotEqual,
			Value:      "12"}).
		Compile()

	if len(deleter.Querys) != 2 {
		t.Error("incorrect querys count. real length: ", len(deleter.Querys))
	}

	if deleter.Querys[0] != "delete from config where state_id = 123 or delegate_state_id != 12;" &&
		deleter.Querys[1] != "delete from config where state_id = 123 or delegate_state_id != 12;" {
		t.Error("incorrect first query: ", deleter.Querys[0])
	}

	if deleter.Querys[0] != "delete from system_recognized_states where state_id = 123 or delegate_state_id != 12;" &&
		deleter.Querys[1] != "delete from system_recognized_states where state_id = 123 or delegate_state_id != 12;" {
		t.Error("incorrect second query: ", deleter.Querys[1])
	}
}

func TestPgBuilderUpdate(t *testing.T) {
	node := createTestNode()
	build := builder.PgBuilder{}

	updater := build.Update(&node.MyStateID, &node.MyWalletID, &node.DelegateStateID, &node.DelegateWalletID).
		Where(builder.Condition{
			Field:      &node.MyStateID,
			Comparator: builder.Equal,
			Value:      node.MyStateID.String()}).
		Or(builder.Condition{
			Field:      &node.DelegateStateID,
			Comparator: builder.NotEqual,
			Value:      "12"}).
		Compile()

	if len(updater.Querys) != 2 {
		t.Error("incorrect querys count. real length: ", len(updater.Querys))
	}

	if updater.Querys[0] != "update config set state_id = 123, wallet_id = 456;" &&
		updater.Querys[1] != "update config set state_id = 123, wallet_id = 456;" {
		t.Error("incorrect first query: ", updater.Querys[0])
	}

	if updater.Querys[0] != "update system_recognized_states set delegate_state_id = 567, delegate_wallet_id = 321 where state_id = 123 or delegate_state_id != 12;" &&
		updater.Querys[1] != "update system_recognized_states set delegate_state_id = 567, delegate_wallet_id = 321 where state_id = 123 or delegate_state_id != 12;" {
		t.Error("incorrect second query: ", updater.Querys[1])
	}
}
