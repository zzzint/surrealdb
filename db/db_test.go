// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"net/http/httptest"

	"github.com/abcum/fibre"
	"github.com/abcum/surreal/cnf"
	"github.com/abcum/surreal/util/data"
)

func setupDB() {

	cnf.Settings = &cnf.Options{}
	cnf.Settings.DB.Path = "memory"
	cnf.Settings.DB.Base = "*"
	cnf.Settings.DB.Proc.Size = 5

	workerCount = 1

	Setup(cnf.Settings)

}

func setupKV() *fibre.Context {

	keep := new(data.Doc)

	auth := new(cnf.Auth)
	auth.Kind = cnf.AuthKV
	auth.Possible.NS = "*"
	auth.Selected.NS = "*"
	auth.Possible.DB = "*"
	auth.Selected.DB = "*"

	req := &fibre.Request{Request: httptest.NewRequest("GET", "/", nil)}
	res := &fibre.Response{}

	ctx := fibre.NewContext(req, res, nil)
	ctx.Set("auth", auth)
	ctx.Set("keep", keep)

	return ctx

}

func setupSC() *fibre.Context {

	keep := new(data.Doc)

	auth := new(cnf.Auth)
	auth.Kind = cnf.AuthSC
	auth.Possible.NS = "*"
	auth.Selected.NS = "*"
	auth.Possible.DB = "*"
	auth.Selected.DB = "*"

	req := &fibre.Request{Request: httptest.NewRequest("GET", "/", nil)}
	res := &fibre.Response{}

	ctx := fibre.NewContext(req, res, nil)
	ctx.Set("auth", auth)
	ctx.Set("keep", keep)

	return ctx

}