package Datastore

// Copyright 2016 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"appengine/datastore"
	"bytes"
	"encoding/gob"
	"fmt"
	"google.golang.org/appengine"

	"github.com/ThePiachu/Go/Log"
)

const FakeBlobstoreBucket = "FakeBlobstore"

type FakeBlobstoreData struct {
	Data []byte `datastore:",noindex"`
}

const MaxDataToStore int = 1000000 //<2**20 to fit into datastore

func PutInFakeBlobstore(c appengine.Context, kind, stringID string, toStore interface{}) error {
	var data bytes.Buffer

	enc := gob.NewEncoder(&data)

	err := enc.Encode(toStore)
	if err != nil {
		Log.Errorf(c, "PutInFakeBlobstore error for key %s:%s - %s", kind, stringID, err)
		return err
	}

	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		for i := 0; ; i++ {
			f := new(FakeBlobstoreData)
			f.Data = data.Next(MaxDataToStore)
			id := makeFakeBlobstoreID(kind, stringID, i)
			_, err := PutInDatastoreSimple(c, FakeBlobstoreBucket, id, f)
			if err != nil {
				Log.Errorf(c, "PutInFakeBlobstore - %v", err)
				return err
			}
			if data.Len() == 0 {
				break
			}
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})

	if err != nil {
		Log.Errorf(c, "PutInFakeBlobstore - %v", err)
		return err
	}

	return nil
}

func GetFromFakeBlobstore(c appengine.Context, kind, stringID string, dst interface{}) error {
	data := bytes.NewBuffer(nil)
	for i := 0; ; i++ {
		id := makeFakeBlobstoreID(kind, stringID, i)
		tmp := new(FakeBlobstoreData)
		err := GetFromDatastoreSimple(c, FakeBlobstoreBucket, id, tmp)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
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
		if len(tmp.Data) < MaxDataToStore {
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
