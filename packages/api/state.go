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

package api

import (
	"fmt"
	"net/http"

	"github.com/EGaaS/go-egaas-mvp/packages/converter"
	"github.com/EGaaS/go-egaas-mvp/packages/model"
	"github.com/EGaaS/go-egaas-mvp/packages/utils/tx"
)

type stateParamResult struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	Conditions string `json:"conditions"`
}

type stateParamListResult struct {
	Count string             `json:"count"`
	List  []stateParamResult `json:"list"`
}

type stateItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Coords string `json:"coords"`
}

type stateListResult struct {
	Count string      `json:"count"`
	List  []stateItem `json:"list"`
}

func getStateParams(w http.ResponseWriter, r *http.Request, data *apiData) error {
	sp := &model.StateParameter{}
	sp.SetTablePrefix(getPrefix(data))
	found, err := sp.GetByName(data.params["name"].(string))
	if !found {
		return errorAPI(w, fmt.Sprintf("State parameter not found %s", data.params["name"].(string)), http.StatusNotFound)
	}
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}
	data.result = &stateParamResult{Name: sp.Name, Value: sp.Value, Conditions: sp.Conditions}
	return nil
}

func txPreNewState(w http.ResponseWriter, r *http.Request, data *apiData) error {
	v := tx.NewState{
		Header:       getSignHeader(`NewState`, data),
		StateName:    data.params[`name`].(string),
		CurrencyName: data.params[`currency`].(string),
	}
	data.result = &forSign{Time: converter.Int64ToStr(v.Time), ForSign: v.ForSign()}
	return nil
}

func txNewState(w http.ResponseWriter, r *http.Request, data *apiData) error {
	header, err := getHeader(`NewState`, data)
	header.StateID = 0
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusBadRequest)
	}

	toSerialize := tx.NewState{
		Header:       header,
		StateName:    data.params[`name`].(string),
		CurrencyName: data.params[`currency`].(string),
	}
	hash, err := sendEmbeddedTx(header.Type, header.UserID, toSerialize)
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}
	data.result = hash
	return nil
}

func txPreNewStateParams(w http.ResponseWriter, r *http.Request, data *apiData) error {
	v := tx.NewStateParameters{
		Header:     getSignHeader(`NewStateParameters`, data),
		Name:       data.params[`name`].(string),
		Value:      data.params[`value`].(string),
		Conditions: data.params[`conditions`].(string),
	}
	data.result = &forSign{Time: converter.Int64ToStr(v.Time), ForSign: v.ForSign()}
	return nil
}

func txPreEditStateParams(w http.ResponseWriter, r *http.Request, data *apiData) error {
	v := tx.EditStateParameters{
		Header:     getSignHeader(`EditStateParameters`, data),
		Name:       data.params[`name`].(string),
		Value:      data.params[`value`].(string),
		Conditions: data.params[`conditions`].(string),
	}
	data.result = &forSign{Time: converter.Int64ToStr(v.Time), ForSign: v.ForSign()}
	return nil
}

func txStateParams(w http.ResponseWriter, r *http.Request, data *apiData) error {
	txName := `NewStateParameters`
	if r.Method == `PUT` {
		txName = `EditStateParameters`
	}
	header, err := getHeader(txName, data)
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusBadRequest)
	}

	var toSerialize interface{}

	if txName == `EditStateParameters` {
		toSerialize = tx.EditStateParameters{
			Header:     header,
			Name:       data.params[`name`].(string),
			Value:      data.params[`value`].(string),
			Conditions: data.params[`conditions`].(string),
		}
	} else {
		toSerialize = tx.NewStateParameters{
			Header:     header,
			Name:       data.params[`name`].(string),
			Value:      data.params[`value`].(string),
			Conditions: data.params[`conditions`].(string),
		}
	}
	hash, err := sendEmbeddedTx(header.Type, header.UserID, toSerialize)
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}
	data.result = hash
	return nil
}

func stateParamsList(w http.ResponseWriter, r *http.Request, data *apiData) error {
	limit := data.params[`limit`].(int64)
	if limit == 0 {
		limit = 25
	} else if limit < 0 {
		limit = -1
	}
	outList := make([]stateParamResult, 0)
	sp := &model.StateParameter{}
	count, err := sp.GetCount(getPrefix(data))
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}

	spList := &model.StateParameter{}
	list, err := spList.GetAllLimitOffset(getPrefix(data), limit, data.params["offset"].(int64))
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}

	for _, val := range list {
		outList = append(outList, stateParamResult{Name: val.Name, Value: val.Value,
			Conditions: val.Conditions})
	}
	data.result = &stateParamListResult{Count: converter.Int64ToStr(count), List: outList}
	return nil
}

func stateList(w http.ResponseWriter, r *http.Request, data *apiData) error {
	limit := data.params[`limit`].(int64)
	if limit == 0 {
		limit = 25
	} else if limit < 0 {
		limit = -1
	}
	ssCount := &model.SystemState{}
	count, err := ssCount.GetCount()
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}
	ssList := &model.SystemState{}
	idata, err := ssList.GetAllLimitOffset(limit, data.params["offset"].(int64))
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}
	outList := make([]stateItem, 0)
	for _, val := range idata {
		if !model.IsNodeState(val.ID, r.Host) {
			continue
		}
		list, err := model.GetAllSystemParametersNameIn([]string{"state_name", "state_flag", "state_coords"})
		if err != nil {
			return errorAPI(w, err.Error(), http.StatusInternalServerError)
		}
		item := stateItem{ID: converter.Int64ToStr(val.ID)}
		for _, val := range list {
			switch val.Name {
			case `state_name`:
				item.Name = val.Value
			case `state_flag`:
				item.Logo = val.Value
			case `state_coords`:
				item.Coords = val.Value
			}
		}
		outList = append(outList, item)
	}
	data.result = &stateListResult{Count: converter.Int64ToStr(count), List: outList}
	return nil
}
