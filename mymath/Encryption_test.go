package mymath

// Copyright 2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	_, err := AESEncrypt("", "")
	if err == nil {
		t.Fail()
		return
	}
	_, err = AESEncrypt("test", "test")
	if err == nil {
		t.Fail()
		return
	}
	resp, err := AESEncrypt("example key 1234", "test text 12345789012345678901234567890")
	if err != nil {
		t.Fail()
		return
	}
	if resp == nil {
		t.Fail()
		return
	}
	if len(resp) == 0 {
		t.Fail()
		return
	}
}

func TestAESDecrypt(t *testing.T) {
	_, err := AESDecrypt("", nil)
	if err == nil {
		t.Fail()
		return
	}
	_, err = AESDecrypt("test", []byte{})
	if err == nil {
		t.Fail()
		return
	}
	resp, err := AESDecrypt("example key 1234", Str2Hex("22277966616d9bc47177bd02603d08c9a67d5380d0fe8cf3b44438dff7b9"))
	if err != nil {
		t.Fail()
	}
	if resp != "some plaintext" {
		t.Fail()
		return
	}
}

func TestAESEncryptAndDecrypt(t *testing.T) {
	resp, err := AESEncrypt("example key 1234", "test text 12345789012345678901234567890")
	if err != nil {
		t.Fail()
	}
	resp2, err := AESDecrypt("example key 1234", resp)
	if err != nil {
		t.Fail()
	}
	if resp2 != "test text 12345789012345678901234567890" {
		t.Fail()
		return
	}
}

func TestAESEncryptAndDecryptCBC(t *testing.T) {
	resp, err := AESEncryptCBC("example key 1234", "test text 1234578901234567890123")
	if err != nil {
		t.Fail()
	}
	resp2, err := AESDecryptCBC("example key 1234", resp)
	if err != nil {
		t.Fail()
	}
	if resp2 != "test text 1234578901234567890123" {
		t.Fail()
		return
	}
}
