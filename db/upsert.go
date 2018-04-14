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
	"fmt"

	"context"

	"github.com/abcum/surreal/sql"
	"github.com/abcum/surreal/util/data"
	"github.com/abcum/surreal/util/keys"
)

func (e *executor) executeUpsert(ctx context.Context, stm *sql.UpsertStatement) ([]interface{}, error) {

	data, err := e.fetch(ctx, stm.Data, nil)
	if err != nil {
		return nil, err
	}

	i := newIterator(e, ctx, stm, false)

	switch data := data.(type) {

	default:
		return nil, fmt.Errorf("Can not execute UPSERT query using value '%v'", data)

	case []interface{}:
		key := &keys.Thing{KV: stm.KV, NS: stm.NS, DB: stm.DB, TB: stm.Into.TB}
		i.processArray(ctx, key, data)

	case map[string]interface{}:
		key := &keys.Thing{KV: stm.KV, NS: stm.NS, DB: stm.DB, TB: stm.Into.TB}
		i.processArray(ctx, key, []interface{}{data})

	}

	return i.Yield(ctx)

}

func (e *executor) fetchUpsert(ctx context.Context, stm *sql.UpsertStatement, doc *data.Doc) (interface{}, error) {

	stm.Echo = sql.AFTER

	if doc != nil {
		vars := data.New()
		vars.Set(doc.Data(), varKeyParent)
		ctx = context.WithValue(ctx, ctxKeySubs, vars)
	}

	out, err := e.executeUpsert(ctx, stm)
	if err != nil {
		return nil, err
	}

	switch len(out) {
	case 1:
		return data.Consume(out).Get(docKeyOne, docKeyId).Data(), nil
	default:
		return data.Consume(out).Get(docKeyAll, docKeyId).Data(), nil
	}

}

func (d *document) runUpsert(ctx context.Context, stm *sql.UpsertStatement) (interface{}, error) {

	var ok bool
	var err error
	var met = _UPDATE

	defer d.close()

	if err = d.setup(); err != nil {
		return nil, err
	}

	if d.val.Exi() == false {
		met = _CREATE
	}

	if ok, err = d.allow(ctx, met); err != nil {
		return nil, err
	} else if ok == false {
		return nil, nil
	}

	if err = d.merge(ctx, met, nil); err != nil {
		return nil, err
	}

	if err = d.storeIndex(); err != nil {
		return nil, err
	}

	if err = d.storeThing(); err != nil {
		return nil, err
	}

	if err = d.table(ctx, met); err != nil {
		return nil, err
	}

	if err = d.event(ctx, met); err != nil {
		return nil, err
	}

	if err = d.lives(ctx, met); err != nil {
		return nil, err
	}

	return d.yield(ctx, stm, stm.Echo)

}