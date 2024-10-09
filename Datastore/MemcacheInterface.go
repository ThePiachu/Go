package Datastore

// Copyright 2012-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/ThePiachu/Go/Log"
	"golang.org/x/net/context"
	"google.golang.org/appengine/v2/capability"
	"google.golang.org/appengine/v2/datastore"
	"google.golang.org/appengine/v2/memcache"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}

func PutInMemcache(c context.Context, key string, toStore interface{}) error {
	if !capability.Enabled(c, "memcache", "*") {
		Log.Errorf(c, "PutInMemcache - Memcache not available.")
		return nil
	}
	var data bytes.Buffer

	enc := gob.NewEncoder(&data)

	err := enc.Encode(toStore)
	if err != nil {
		Log.Errorf(c, "PutInMemcache error for key %s - %s", key, err)
		return err
	}

	item := &memcache.Item{
		Key:   key,
		Value: data.Bytes(),
	}
	if err := memcache.Set(c, item); err != nil {
		Log.Errorf(c, "PutInMemcache - %s", err)
	}
	return err
}

func GetFromMemcache(c context.Context, key string, dst interface{}) interface{} {
	if !capability.Enabled(c, "memcache", "*") {
		Log.Errorf(c, "GetFromMemcache - Memcache not available.")
		return nil
	}

	item, err := memcache.Get(c, key)
	if err != nil && err != memcache.ErrCacheMiss {
		Log.Errorf(c, "GetFromMemcache - %s", err)
		return nil
	} else if err == memcache.ErrCacheMiss {
		return nil
	}

	dec := gob.NewDecoder(bytes.NewBuffer(item.Value))
	err = dec.Decode(dst)
	if err != nil {
		Log.Errorf(c, "GetFromMemcache - %s", err)
		return nil
	}

	return dst
}

func PutInDatastoreSimpleAndMemcache(c context.Context, kind, stringID, memcacheID string, variable interface{}) (*datastore.Key, error) {
	if !capability.Enabled(c, "datastore_v3", "*") {
		Log.Errorf(c, "PutInDatastoreSimpleAndMemcache - Datastore not available.")
		return nil, errors.New("Datastore not available")
	}

	key, err := PutInDatastoreSimple(c, kind, stringID, variable)
	if err != nil {
		Log.Errorf(c, "PutInDatastoreSimpleAndMemcache - %s", err)
		return nil, err
	}

	if capability.Enabled(c, "memcache", "*") {
		PutInMemcache(c, memcacheID, variable)
	}

	return key, nil
}

func GetFromDatastoreSimpleOrMemcache(c context.Context, kind, stringID, memcacheID string, dst interface{}) error {
	if capability.Enabled(c, "memcache", "*") {
		answer := GetFromMemcache(c, memcacheID, dst)
		if answer != nil {
			dst = answer
			return nil
		}
	}
	if !capability.Enabled(c, "datastore_v3", "*") {
		Log.Errorf(c, "GetFromDatastoreOrMemcache - Datastore not available.")
		return errors.New("Datastore not available")
	}

	err := GetFromDatastoreSimple(c, kind, stringID, dst)
	if err != nil {
		Log.Infof(c, "Problem getting entity %v of kind %v from datastore: %s", stringID, kind, err)
		return err
	}

	if capability.Enabled(c, "memcache", "*") {
		PutInMemcache(c, memcacheID, dst)
	}

	return nil
}

func IsVariableInDatastoreSimpleOrMemcache(c context.Context, kind, stringID, memcacheID string, dst interface{}) bool {
	_, err := memcache.Get(c, memcacheID)
	if err == nil {
		return true
	}
	return IsVariableInDatastoreSimple(c, kind, stringID, dst)
}

func DeleteFromMemcache(c context.Context, memcacheID string) {
	memcache.Delete(c, memcacheID)
}

func DeleteFromDatastoreSimpleAndMemcache(c context.Context, kind, stringID, memcacheID string) error {
	DeleteFromMemcache(c, memcacheID)
	return DeleteFromDatastoreSimple(c, kind, stringID)
}

func FlushMemcache(c context.Context) error {
	return memcache.Flush(c)
}

func ClearNamespaceAndMemcache(c context.Context, kind string) error {
	err := ClearNamespace(c, kind)
	if err != nil {
		Log.Errorf(c, "ClearNamespaceAndMemcache - %v", err)
		return err
	}
	err = FlushMemcache(c)
	if err != nil {
		Log.Errorf(c, "ClearNamespaceAndMemcache - %v", err)
		return err
	}
	return nil
}

func TestMemcache(c context.Context) {
	type TMP struct {
		A string
		B int
		C float64
	}
	tmp := new(TMP)
	tmp.A = "Hello"
	tmp.B = 123
	tmp.C = 12.3
	Log.Infof(c, "TestMemcache")
	key, err := PutInDatastoreSimpleAndMemcache(c, "test", "test", "test", tmp)
	Log.Infof(c, "%v, %v", key, err)
	Log.Infof(c, "%v", GetFromDatastoreSimpleOrMemcache(c, "test", "test", "test", tmp))
	Log.Infof(c, "%v", tmp)
}
