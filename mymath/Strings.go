package mymath

// Copyright 2011-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"sort"
)

type Alphabetic []string

func (a Alphabetic) Len() int {
	return len(a)
}
func (a Alphabetic) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a Alphabetic) Less(i, j int) bool {
	var max int
	if len(a[i]) < len(a[j]) {
		max = len(a[i])
	} else {
		max = len(a[j])
	}
	for k := 0; k < max; k++ {
		if a[i][k] != a[j][k] {
			return a[i][k] < a[j][k]
		}
	}
	return len(a[i]) < len(a[j])
}

func SortStringList(list []string) {
	sort.Sort(Alphabetic(list))
}
