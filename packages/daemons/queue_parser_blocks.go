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

package daemons

import (
	"github.com/EGaaS/go-egaas-mvp/packages/consts"
	"github.com/EGaaS/go-egaas-mvp/packages/model"
	"github.com/EGaaS/go-egaas-mvp/packages/parser"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"

	"context"
)

/* Take the block. If the block has the best hash, then look for the block where the fork started
 * If the fork begins less then variables->rollback_blocks blocks ago, than
 *  - get the whole chain of blocks
 *  - roll back the frontal data from our blocks
 *  - insert the frontal data from a new chain
 *  - if there is no error, then roll back our data from the blocks
 *  - and insert new data
 *  - if there are errors, then roll back to the former data
 * if the fork was long ago then do not touch anything and leave the script blocks_collection.php
 * the limitation variables->rollback_blocks is needed for the protection against the false blocks
 *
 * */

// QueueParserBlocks parses blocks from the queue
func QueueParserBlocks(d *daemon, ctx context.Context) error {

	locked, err := DbLock(ctx, d.goRoutineName)
	if !locked || err != nil {
		return err
	}
	defer DbUnlock(d.goRoutineName)

	infoBlock := &model.InfoBlock{}
	err = infoBlock.GetInfoBlock()
	if err != nil {
		return err
	}
	queueBlock := &model.QueueBlock{}
	err = queueBlock.GetQueueBlock()
	if err != nil && err != model.RecordNotFound {
		return err
	} else if err == model.RecordNotFound {
		return nil
	}
	// check if the block gets in the rollback_blocks_1 limit
	if queueBlock.BlockID > infoBlock.BlockID+consts.RB_BLOCKS_1 {
		queueBlock.Delete()
		return utils.ErrInfo("rollback_blocks_1")
	}

	// is it old block in queue ?
	if queueBlock.BlockID <= infoBlock.BlockID {
		queueBlock.Delete()
		return utils.ErrInfo("old block")
	}

	// download blocks for check
	fullNode := &model.FullNode{}

	err = fullNode.FindNodeByID(queueBlock.FullNodeID)
	if err != nil {
		queueBlock.Delete()
		return utils.ErrInfo(err)
	}

	blockID := queueBlock.BlockID

	p := new(parser.Parser)
	p.GoroutineName = d.goRoutineName

	host := GetHostPort(fullNode.Host)
	err = p.GetBlocks(blockID, host, "rollback_blocks_1", d.goRoutineName, 7)
	if err != nil {
		log.Error("v", err)
		queueBlock.Delete()
		return utils.ErrInfo(err)
	}
	return nil
}
