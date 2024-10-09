package Datastore_test

// Copyright 2016 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"crypto/rand"
	"fmt"
	"google.golang.org/appengine/v2/aetest"
	"testing"

	. "github.com/ThePiachu/Go/Datastore"
)

func TestFakeBlobstore(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	type temp struct {
		ToStore string
		Blarg   int
	}
	b := new(temp)
	b.Blarg = 100
	data := make([]byte, MaxDataToStore*3)

	_, err = rand.Read(data)
	if err != nil {
		t.Errorf("%v", err)
	}
	b.ToStore = fmt.Sprintf("%x", data)

	id := "id"
	kind := "kind"

	err = PutInFakeBlobstore(c, kind, id, b)
	if err != nil {
		t.Errorf("%v", err)
	}

	b2 := new(temp)

	err = GetFromFakeBlobstore(c, kind, id, b2)
	if err != nil {
		t.Errorf("%v", err)
	}

	if b2.Blarg != b.Blarg {
		t.Errorf("Ints don't match")
	}

	if b2.ToStore != b.ToStore {
		t.Errorf("Strings don't match")
	}
}
