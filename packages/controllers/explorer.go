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

package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"

	"github.com/EGaaS/go-egaas-mvp/packages/static"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
)

type Public struct {
	Content template.HTML
}

func Explorer(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Content Recovered", fmt.Sprintf("%s: %s", e, debug.Stack()))
		}
	}()
	w.Header().Set("Content-type", "text/html")
	c := new(Controller)
	c.r = r
	c.w = w
	c.DCDB = utils.DB
	c.ContentInc = true
	c.dbInit = true

	r.ParseForm()
	pageName := `block_explorer`
	c.Parameters, _ = c.GetParameters()

	/*	lang := GetLang(w, r, c.Parameters)
		c.Lang = globalLangReadOnly[lang]
		c.LangInt = int64(lang)
		if lang == 42 {
			c.TimeFormat = "2006-01-02 15:04:05"
		} else {
			c.TimeFormat = "2006-02-01 15:04:05"
		}

		c.Periods = map[int64]string{86400: "1 " + c.Lang["day"], 604800: "1 " + c.Lang["week"], 31536000: "1 " + c.Lang["year"], 2592000: "1 " + c.Lang["month"], 1209600: "2 " + c.Lang["weeks"]}
	*/
	c.Data = &CommonPage{
		Address:      c.SessAddress,
		WalletId:     c.SessWalletId,
		CitizenId:    c.SessCitizenId,
		StateId:      c.SessStateId,
		StateName:    ``,
		CountSignArr: []int{0}, // !!! Добавить вычисление
	}
	content := CallPage(c, pageName)
	funcMap := template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	data, err := static.Asset("static/public.html")
	t := template.New("template").Funcs(funcMap)
	t, err = t.Parse(string(data))
	if err != nil {
		log.Error(err.Error())
	}

	b := new(bytes.Buffer)
	err = t.Execute(b, &Public{Content: template.HTML(content)})
	if err != nil {
		log.Error(err.Error())
	}
	w.Write(b.Bytes())
}
