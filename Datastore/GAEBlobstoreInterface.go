package Datastore

// Copyright 2013-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"encoding/gob"
	"github.com/ThePiachu/Go/Log"
	"golang.org/x/net/context"
	appengine "google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/blobstore"
	"io/ioutil"
)

/*
//TODO: fix
func PutInBlobstore(c context.Context, toStore interface{}) (appengine.BlobKey, error) {
	//TODO: check capabilities
	var k appengine.BlobKey
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(toStore)
	if err != nil {
		Log.Errorf(c, "PutInBlobstore - %s", err)
		return k, err
	}

	w, err := blobstore.Create(c, "application/octet-stream")
	if err != nil {
		Log.Errorf(c, "PutInBlobstore - %s", err)
		return k, err
	}
	_, err = w.Write(data.Bytes())
	if err != nil {
		Log.Errorf(c, "PutInBlobstore - %s", err)
		return k, err
	}
	err = w.Close()
	if err != nil {
		Log.Errorf(c, "PutInBlobstore - %s", err)
		return k, err
	}
	k, err = w.Key()
	if err != nil {
		Log.Errorf(c, "PutInBlobstore - %s", err)
	}
	return k, err
}
*/
func GetFromBlobstore(c context.Context, blobkey appengine.BlobKey, dst interface{}) (interface{}, error) {
	//TODO: check capabilities

	reader := blobstore.NewReader(c, blobkey)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		Log.Errorf(c, "GetFromBlobstore - %s", err)
		return nil, err
	}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err = dec.Decode(dst)
	if err != nil {
		Log.Errorf(c, "GetFromBlobstore - %s", err)
		return nil, err
	}
	return dst, nil
}

func DeleteFromBlobstore(c context.Context, blobkey appengine.BlobKey) error {
	return blobstore.Delete(c, blobkey)
}
