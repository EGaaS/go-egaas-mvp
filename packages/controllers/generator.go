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

const nGenerator = `generator`

type generatorPage struct {
	Data   *CommonPage
	Params map[string]string
}

func init() {
	newPage(nGenerator)
}

// Generator is a control for creating a new citizen
func (c *Controller) Generator() (string, error) {
	pageData := generatorPage{Data: c.Data, Params: make(map[string]string)}
	for key, val := range c.r.Form {
		pageData.Params[key] = val[0]
	}
	return proceedTemplate(c, nGenerator, &pageData)
}
