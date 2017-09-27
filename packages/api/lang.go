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

type langResult struct {
	Name  string `json:"name"`
	Trans string `json:"trans"`
}

type langListResult struct {
	Count string       `json:"count"`
	List  []langResult `json:"list"`
}

func getLang(w http.ResponseWriter, r *http.Request, data *apiData) error {
	prefix := getPrefix(data)
	lang := &model.Language{}
	lang.SetTablePrefix(prefix)
	found, err := lang.Get(data.params["name"].(string))
	if !found {
		return errorAPI(w, fmt.Sprintf("Lang not found %s", data.params["name"].(string)), http.StatusNotFound)
	}
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}
	data.result = &langResult{Name: lang.Name, Trans: lang.Res}
	return nil
}

func txPreNewLang(w http.ResponseWriter, r *http.Request, data *apiData) error {
	v := tx.EditNewLang{
		Header: getSignHeader(`NewLang`, data),
		Name:   data.params[`name`].(string),
		Trans:  data.params[`trans`].(string),
	}
	data.result = &forSign{Time: converter.Int64ToStr(v.Time), ForSign: v.ForSign()}
	return nil
}

func txPreEditLang(w http.ResponseWriter, r *http.Request, data *apiData) error {
	v := tx.EditNewLang{
		Header: getSignHeader(`EditLang`, data),
		Name:   data.params[`name`].(string),
		Trans:  data.params[`trans`].(string),
	}
	data.result = &forSign{Time: converter.Int64ToStr(v.Time), ForSign: v.ForSign()}
	return nil
}

func txLang(w http.ResponseWriter, r *http.Request, data *apiData) error {
	txName := `NewLang`
	if r.Method == `PUT` {
		txName = `EditLang`
	}
	header, err := getHeader(txName, data)
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusBadRequest)
	}

	toSerialize := tx.EditNewLang{
		Header: header,
		Name:   data.params[`name`].(string),
		Trans:  data.params[`trans`].(string),
	}
	hash, err := sendEmbeddedTx(header.Type, header.UserID, toSerialize)
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}
	data.result = hash
	return nil
}

func langList(w http.ResponseWriter, r *http.Request, data *apiData) error {
	limit := data.params[`limit`].(int64)
	if limit == 0 {
		limit = 25
	} else if limit < 0 {
		limit = -1
	}
	outList := make([]langResult, 0)
	langCount := &model.Language{}
	count, err := langCount.GetCount(getPrefix(data))
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}

	langList := &model.Language{}
	list, err := langList.GetAllLimitOffset(getPrefix(data), limit, data.params["offset"].(int64))
	if err != nil {
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}

	for _, val := range list {
		outList = append(outList, langResult{Name: val.Name, Trans: val.Res})
	}
	data.result = &langListResult{Count: converter.Int64ToStr(count), List: outList}
	return nil
}
