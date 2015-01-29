package mymath

// Copyright 2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"testing"
)

func TestRoundFloatUpToXSignificantDigits(t *testing.T) {
	if RoundFloatUpToXSignificantDigits(1, 1) != 2 {
		t.Errorf("Failed 1, 1")
		t.Fail()
	}
	if RoundFloatUpToXSignificantDigits(123, 1) != 200 {
		t.Errorf("Failed 123, 1")
		t.Fail()
	}
	if RoundFloatUpToXSignificantDigits(0.987654, 3) != 0.988 {
		t.Errorf("Failed 0.987654, 3 - %v", RoundFloatUpToXSignificantDigits(0.987654, 3))
		t.Fail()
	}
	if RoundFloatUpToXSignificantDigits(1.00000001, 1) != 2 {
		t.Errorf("Failed 1.00000001, 1")
		t.Fail()
	}
	if RoundFloatUpToXSignificantDigits(1.00000001, 3) != 1.01 {
		t.Errorf("Failed 1.00000001, 3")
		t.Fail()
	}
	if RoundFloatUpToXSignificantDigits(1.9999999, 3) != 2 {
		t.Errorf("Failed 1.9999999, 3")
		t.Fail()
	}
}
