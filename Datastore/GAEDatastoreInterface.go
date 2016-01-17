package Datastore

// Copyright 2012-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"appengine"
	"appengine/datastore"
	"github.com/ThePiachu/Go/Log"
)

func PutInDatastoreFull(c appengine.Context, kind, stringID string, intID int64, parent *datastore.Key, variable interface{}) (*datastore.Key, error) {
	k := datastore.NewKey(c, kind, stringID, intID, parent)
	key, err := datastore.Put(c, k, variable)
	return key, err
}

func PutInDatastoreSimple(c appengine.Context, kind, stringID string, variable interface{}) (*datastore.Key, error) {
	return PutInDatastoreFull(c, kind, stringID, 0, nil, variable)
}

func PutInDatastoreSimpleWithParent(c appengine.Context, kind, stringID string, parent *datastore.Key, variable interface{}) (*datastore.Key, error) {
	return PutInDatastoreFull(c, kind, stringID, 0, parent, variable)
}

func PutInDatastore(c appengine.Context, kind string, variable interface{}) (*datastore.Key, error) {
	return PutInDatastoreFull(c, kind, "", 0, nil, variable)
}

func GetFromDatastoreFull(c appengine.Context, kind, stringID string, intID int64, parent *datastore.Key, dst interface{}) error {
	k := datastore.NewKey(c, kind, stringID, intID, parent)
	return datastore.Get(c, k, dst)
}

func GetFromDatastoreSimple(c appengine.Context, kind, stringID string, dst interface{}) error {
	return GetFromDatastoreFull(c, kind, stringID, 0, nil, dst)
}

func GetFromDatastoreSimpleWithParent(c appengine.Context, kind, stringID string, parent *datastore.Key, dst interface{}) error {
	return GetFromDatastoreFull(c, kind, stringID, 0, parent, dst)
}

// A function that either loads a variable from datastore, or if it is not present, sets it and then loads it
func GetFromDatastoreOrSetDefaultFull(c appengine.Context, kind, stringID string, intID int64, parent *datastore.Key, dst interface{}, def interface{}) error {

	key := datastore.NewKey(c, kind, stringID, intID, parent)
	if err := datastore.Get(c, key, dst); err != nil {
		if err == datastore.ErrNoSuchEntity {
			_, err2 := datastore.Put(c, key, def)

			if err2 != nil {
				return err2
			} else {
				if err3 := datastore.Get(c, key, dst); err3 != nil {
					return err3
				}
			}

		} else {
			return err
		}
	}
	return nil
}

func GetFromDatastoreOrSetDefaultSimple(c appengine.Context, kind, stringID string, dst interface{}, def interface{}) error {
	return GetFromDatastoreOrSetDefaultFull(c, kind, stringID, 0, nil, dst, def)
}

func IsVariableInDatastoreSimple(c appengine.Context, kind, stringID string, dst interface{}) bool {
	err := GetFromDatastoreSimple(c, kind, stringID, dst)
	if err == nil {
		return true
	}
	if err == datastore.ErrNoSuchEntity {
		return false
	}
	Log.Errorf(c, "IsVariableInDatastoreSimple - %s", err)
	return false
}

func QueryGetFirstBy(c appengine.Context, kind string, order string, limit int, dst interface{}) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind).Order(order).Limit(limit)
	return q.GetAll(c, dst)
}

func QueryGetFirstKeysBy(c appengine.Context, kind string, order string, limit int, dst interface{}) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind).Order(order).Limit(limit).KeysOnly()
	return q.GetAll(c, dst)
}
func QueryGetAllWithFilter(c appengine.Context, kind string, filterStr string, filterValue interface{}, dst interface{}) ([]*datastore.Key, error) {
	return QueryGetAllWithFilterAndLimit(c, kind, filterStr, filterValue, -1, dst)
}

func QueryGetAllWithFilterAndLimit(c appengine.Context, kind string, filterStr string, filterValue interface{}, limit int, dst interface{}) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind).Filter(filterStr, filterValue).Limit(limit)
	return q.GetAll(c, dst)
}

func QueryGetAllWithLimit(c appengine.Context, kind string, limit int, dst interface{}) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind).Limit(limit)
	return q.GetAll(c, dst)
}

func QueryGetAll(c appengine.Context, kind string, dst interface{}) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind)
	return q.GetAll(c, dst)
}

func QueryGetAllKeysWithFilterAndOrder(c appengine.Context, kind string, filterStr string, filterValue interface{}, orderStr string, dst interface{}) ([]*datastore.Key, error) {
	return QueryGetAllKeysWithFilterLimitAndOrder(c, kind, filterStr, filterValue, -1, orderStr, dst)
}

func QueryGetAllKeysWithFilterLimitAndOrder(c appengine.Context, kind string, filterStr string, filterValue interface{}, limit int, orderStr string, dst interface{}) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind).Filter(filterStr, filterValue).Limit(limit).Order(orderStr).KeysOnly()
	return q.GetAll(c, dst)
}

func QueryGetAllKeysWithFilterLimitOffsetAndOrder(c appengine.Context, kind string, filterStr string, filterValue interface{}, limit int, offset int, orderStr string, dst interface{}) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind).Filter(filterStr, filterValue).Limit(limit).Offset(offset).Order(orderStr).KeysOnly()
	return q.GetAll(c, dst)
}

func QueryGetAllKeysWithFilter(c appengine.Context, kind string, filterStr string, filterValue interface{}, dst interface{}) []*datastore.Key {
	return QueryGetAllKeysWithFilterAndLimit(c, kind, filterStr, filterValue, -1, dst)
}

func QueryGetAllKeysWithFilterAndLimit(c appengine.Context, kind string, filterStr string, filterValue interface{}, limit int, dst interface{}) []*datastore.Key {
	q := datastore.NewQuery(kind).Filter(filterStr, filterValue).Limit(limit).KeysOnly()
	keys, err := q.GetAll(c, dst)
	if err != nil {
		panic(err)
	}
	return keys
}

func QueryGetAllKeys(c appengine.Context, kind string, dst interface{}) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind).KeysOnly()
	return q.GetAll(c, dst)
}

func CountQueryWithFilter(c appengine.Context, kind string, filterStr string, filterValue interface{}) int {
	q := datastore.NewQuery(kind).Filter(filterStr, filterValue)
	count, err := q.Count(c)
	if err != nil {
		Log.Errorf(c, "CountQueryWithFilter - %s", err)
		return -1
	}
	return count
}

func ClearNamespace(c appengine.Context, kind string) error {
	q := datastore.NewQuery(kind)
	q = q.KeysOnly()

	keys, err := q.GetAll(c, nil)

	if err != nil {
		Log.Errorf(c, "Clear Namespace - %v", err)
		return err
	}

	for {
		toDelete := keys[:]
		if len(keys) > 500 {
			toDelete = keys[:500]
			keys = keys[500:]
		}
		err = datastore.DeleteMulti(c, toDelete)
		if err != nil {
			Log.Errorf(c, "ClearNamespace - %v", err)
			return err
		}
		if len(keys) < 500 {
			break
		}
	}

	return nil
}

func DeleteFromDatastoreFull(c appengine.Context, kind, stringID string, intID int64, parent *datastore.Key) error {
	k := datastore.NewKey(c, kind, stringID, intID, parent)
	return datastore.Delete(c, k)
}

func DeleteFromDatastoreSimple(c appengine.Context, kind, stringID string) error {
	return DeleteFromDatastoreFull(c, kind, stringID, 0, nil)
}

func KeysToStringIDs(keys []*datastore.Key) []string {
	answer := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		answer[i] = keys[i].StringID()
	}
	return answer
}
