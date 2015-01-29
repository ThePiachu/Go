package mymath

// Copyright 2011-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	//rand2 "math/rand"
	"crypto/rand"
	"math"
	"math/big"
)

func RandFloat64Between(min float64, max float64) float64 {
	number := RandFloat64()
	number = number * (max - min)
	number = number + min
	return number
}

func RandFloat64() float64 {
	return float64(RandInt64()) / float64(MaxInt64())
}

func Randuint64() []byte {
	uint64max := big.NewInt(1)
	uint64max.Lsh(uint64max, 64)

	randnum, _ := rand.Int(rand.Reader, uint64max)

	random := randnum.Bytes()
	answer := make([]byte, 8)

	for i := 0; i < len(random); i++ {
		answer[i] = random[i]
	}

	return answer
}

func Randuint64Rev() []byte {
	uint64max := big.NewInt(1)
	uint64max.Lsh(uint64max, 64)

	randnum, _ := rand.Int(rand.Reader, uint64max)

	random := randnum.Bytes()
	answer := make([]byte, 8)

	for i := 0; i < len(random); i++ {
		answer[len(answer)-1-i] = random[i]
	}

	return answer
}

func RandInt64() int64 {
	int64max := big.NewInt(math.MaxInt64)
	randnum, _ := rand.Int(rand.Reader, int64max)
	return randnum.Int64()
}

func RandInt64Between(min int64, max int64) int64 {
	if min > max {
		tmp := min
		min = max
		max = tmp
	}
	if min == max {
		return max
	}
	int64max := big.NewInt(max - min)
	int64min := big.NewInt(min)
	randnum, _ := rand.Int(rand.Reader, int64max)
	randnum = randnum.Add(randnum, int64min)
	return randnum.Int64()
}

func RandIntBetween(min int, max int) int {
	return int(RandInt64Between(int64(min), int64(max)))
}

func MaxInt() int {
	return math.MaxInt32
}

func MaxInt64() int64 {
	return math.MaxInt64
}

func RandInt() int {
	return int(RandInt64())
}

func AbsInt64(i int64) int64 {
	if i > 0 {
		return i
	}
	return -i
}

func AbsInt(i int) int {
	if i > 0 {
		return i
	}
	return -i
}

func RandIntPositive() int {
	return AbsInt(RandInt())
}

func RandInt64Positive() int64 {
	return AbsInt64(RandInt64())
}

func RandomHex(length int) []byte {
	if length <= 0 {
		return nil
	}
	answer := make([]byte, length)
	/*
		for i:=0;i<length;i++{
			answer[i]=byte(rand2.Intn(0xFF))
		}
		return answer*/
	_, err := rand.Read(answer)
	if err != nil {
		return nil
	}
	return answer
}

func RandomHexString(length int) string {
	hex := RandomHex(length/2 + 1)
	hexStr := Hex2Str(hex)
	return hexStr[0:length]
}

func RandomAlphabetString(length int, alphabet string) string {
	answer := ""
	alphabetLength := len(alphabet)
	if alphabetLength == 0 {
		return ""
	}
	for i := 0; i < length; i++ {
		letter := RandIntBetween(0, alphabetLength)
		answer = answer + string(alphabet[letter])
	}
	return answer
}
