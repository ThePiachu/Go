package Datastore

// Copyright 2016 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/ThePiachu/Go/Log"
)

const FakeBlobstoreBucket = "FakeBlobstore"

type FakeBlobstoreData struct {
	Data []byte `datastore:",noindex"`
}

const MaxDataToStore int = 1000000 //<2**20 to fit into datastore

func PutInFakeBlobstore(c appengine.Context, kind, stringID string, toStore interface{}) (error) {
	var data bytes.Buffer

	enc := gob.NewEncoder(&data)

	err := enc.Encode(toStore)
	if err != nil {
		Log.Errorf(c, "PutInFakeBlobstore error for key %s:%s - %s", kind, stringID, err)
		return err
	}

	toSave:=make([]*FakeBlobstoreData, 0, data.Len()/MaxDataToStore+1)
	keys:=make([]*datastore.Key, 0, data.Len()/MaxDataToStore+1)

	for i:=0;;i++ {
		f:=new(FakeBlobstoreData)
		f.Data = data.Next(MaxDataToStore)
		id:=makeFakeBlobstoreID(kind, stringID, i)
		
		k := datastore.NewKey(c, FakeBlobstoreBucket, id, 0, nil)
		toSave = append(toSave, f)
		keys = append(keys, k)

		if data.Len()==0 {
			break
		}
	}

	_, err = datastore.PutMulti(c, keys, toSave)

	if err!=nil {
		Log.Errorf(c, "PutInFakeBlobstore - %v", err)
		return err
	}

	return nil
}

func GetFromFakeBlobstore(c appengine.Context, kind, stringID string, dst interface{}) (error) {
	data:=bytes.NewBuffer(nil)
	for i:=0;;i++ {
		id:=makeFakeBlobstoreID(kind, stringID, i)
		tmp:=new(FakeBlobstoreData)
		err:=GetFromDatastoreSimple(c, FakeBlobstoreBucket, id, tmp)
		if err!=nil {
			if err==datastore.ErrNoSuchEntity {
				break
			}
			Log.Errorf(c, "GetFromFakeBlobstore - %v", err)
			return err
		}
		_, err = data.Write(tmp.Data)
		if err != nil {
			Log.Errorf(c, "GetFromFakeBlobstore - %s", err)
			return err
		}
		if len(tmp.Data)<MaxDataToStore {
			break
		}
	}

	dec := gob.NewDecoder(data)
	err := dec.Decode(dst)
	if err != nil {
		Log.Errorf(c, "GetFromFakeBlobstore - %s", err)
		return err
	}

	return nil
}

func makeFakeBlobstoreID(kind, stringID string, i int) string {
	return fmt.Sprintf("%s:%s:%d", kind, stringID, i)
}