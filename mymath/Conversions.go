package mymath

// Copyright 2011-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"speter.net/go/exp/math/dec/inf"
	"strconv"
	"strings"
)

func TrimFloatToXSignificantDigits(f float64, digits int) float64 {
	var tmp string = "%%.%dg"
	tmp = fmt.Sprintf(tmp, digits)
	tmp = fmt.Sprintf(tmp, f)
	return Str2Float(tmp)
}

func RoundFloatUpToXSignificantDigits(f float64, digits int) float64 {
	power := int(Float2Int(math.Floor(math.Log10(f)))) - digits + 1
	f += math.Pow10(power) * 0.5
	var tmp string = "%%.%dg"
	tmp = fmt.Sprintf(tmp, digits)
	tmp = fmt.Sprintf(tmp, f)
	result := Str2Float(tmp)
	return result
}

func TrimFloatToXDecimalDigits(f float64, digits int) float64 {
	var tmp string = "%%.%df"
	tmp = fmt.Sprintf(tmp, digits)
	tmp = fmt.Sprintf(tmp, f)
	return Str2Float(tmp)
}

func Base642Hex(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

func Hex2Base64(b []byte) string {
	s := base64.StdEncoding.EncodeToString(b)
	return s
}

func Str2Base64(s string) string {
	b := []byte(s)
	return Hex2Base64(b)
}

func Str2Int(s string) int {
	return String2Int(s)
}

func String2Int(s string) int {
	return int(String2Int64(s))
}

func Str2Int64(s string) int64 {
	return String2Int64(s)
}

func String2Int64(s string) int64 {
	answer, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Printf("Str2Int64 - %s", err)
		return 0
	}
	return answer
}

func String2Dec(s string) *inf.Dec {
	answer, ok := new(inf.Dec).SetString(s)
	if !ok {
		log.Printf("String2Dec failed for %s", s)
		return nil
	}
	return answer
}

func Str2Dec(s string) *inf.Dec {
	return String2Dec(s)
}

func Float64ToDec(f float64) *inf.Dec {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	answer, _ := new(inf.Dec).SetString(s)
	return answer
}

func HexString2Int64(s string) int64 {
	answer, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		log.Printf("HexString2Int64 - %s", err)
		return 0
	}
	return answer
}

func String2Float(s string) float64 {
	answer, _ := strconv.ParseFloat(s, 64)
	return answer
}

func Str2Float(s string) float64 {
	return String2Float(s)
}

func Float2CurrencyString(f float64) string {
	return Float2CurrencyStr(f)
}

func Float2CurrencyStr(f float64) string {
	number:=Float2Str(f)
	parts:=SplitStrings(number, ".")
	if len(parts) == 0 {
		return number
	}
	whole:=parts[0]
	decimal:=""
	if len(parts)>1 {
		decimal = parts[1]
	}
	length:=len(whole)
	answer:=""
	for k:=range(whole) {
		answer = answer+string(whole[k])
		if ((length-k-1)%3==0) {
			if (length-k-1)!=0{
				answer = answer+"'"
			}
		}
	}
	if decimal!="" {
		answer = answer+"."+decimal
	}
	return answer
}

func Float2Str(f float64) string {
	return Float642String(f)
}

func Float2String(f float64) string {
	return Float642String(f)
}

func Float642String(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Float642Str(f float64) string {
	return Float642String(f)
}

func Float642Int64(f float64) int64 {
	return int64(f)
}

func Float2Int(f float64) int64 {
	return Float642Int64(f)
}

func Int642Float64(i int64) float64 {
	return float64(i)
}
func Int2Float(i int64) float64 {
	return Int642Float64(i)
}

func Int642String(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Int642HexString(i int64) string {
	return strconv.FormatInt(i, 16)
}

func Int642Str(i int64) string {
	return Int642String(i)
}

func Int2Str(i int) string {
	return Int642String(int64(i))
}

func Int2String(i int) string {
	return Int642String(int64(i))
}
func Int2Hex(i int) []byte {
	return Int642Hex(int64(i))
}

func Str2Uint32(s string) uint32 {
	return Hex2Uint32(Str2Hex(s))
}

func MapStringInt64ToString(m map[string]int64) string {
	answer := "{"
	for k, v := range m {
		answer += "\"" + k + "\":" + Int642Str(v) + ","
	}
	answer = answer[:len(answer)-1]
	answer += "}"
	return answer
}

func Hex2Int(b []byte) int {
	var answer int
	answer = 0
	if len(b) > 0 {
		maxBytes := len(b)
		for i := 0; i < maxBytes; i++ {
			answer *= 256
			answer += int(b[i])
		}
	}
	return answer
}

func Hex2Uint64(b []byte) uint64 {
	var answer uint64
	answer = 0
	if len(b) > 0 {
		maxBytes := 8
		if len(b) < 8 {
			maxBytes = len(b)
		}
		for i := 0; i < maxBytes; i++ {
			answer *= 256
			answer += uint64(b[i])
		}
	}
	return answer
}

func Hex2Uint32(b []byte) uint32 {
	var answer uint32
	answer = 0
	if len(b) > 0 {
		maxBytes := 4
		if len(b) < 4 {
			maxBytes = len(b)
		}
		for i := 0; i < maxBytes; i++ {
			answer *= 256
			answer += uint32(b[i])
		}
	}
	return answer
}

func HexRev2Uint64(b []byte) uint64 {
	return Hex2Uint64(Rev(b))
}

func HexRev2Uint32(b []byte) uint32 {
	return Hex2Uint32(Rev(b))
}

func PadHexToLength(hex []byte, length int) []byte {
	if len(hex) > length {
		return hex
	}
	answer := make([]byte, length)
	for i := 0; i < len(hex); i++ {
		answer[len(answer)-1-i] = hex[len(hex)-1-i]
	}
	return answer
}

func Big2HexWithMinimumLength(b *big.Int, length int) []byte {
	return PadHexToLength(b.Bytes(), length)
}

//TODO: test
func Big2Hex(b *big.Int) []byte {
	return b.Bytes()
}

//TODO: test
func Big2HexRev(b *big.Int) []byte {
	return Rev(Big2Hex(b))
}

func String2Hex32(s string) [32]byte {
	var answer [32]byte
	if len(s) == 64 {
		copy(answer[:], String2Hex(s))
	}
	return answer
}

func String2Hex(s string) []byte {
	if len(s) == 0 {
		return []byte{0}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	answer, _ := hex.DecodeString(s)
	return answer
}

func Str2Hex(s string) []byte {
	return String2Hex(s)
}

func String2HexRev(s string) []byte {
	answer, _ := hex.DecodeString(s)
	return Rev(answer)
}

func Str2HexRev(s string) []byte {
	return Rev(String2Hex(s))
}

//TODO: test
func String2BigBase(s string, base int) *big.Int {
	b := big.NewInt(0)
	b.SetString(s, base)
	return b
}

//TODO: test
func Str2BigBase(s string, base int) *big.Int {
	return String2BigBase(s, base)
}

//TODO: test
func String2Big(s string) *big.Int {
	return String2BigBase(s, 16)
}

//TODO: test
func Str2Big(s string) *big.Int {
	return String2BigBase(s, 16)
}

//TODO: test
func ASCII2Hex(s string) []byte {
	return []byte(s)
}
func ASCII2HexRev(s string) []byte {
	return Rev(ASCII2Hex(s))
}

//TODO: test
func Hex2ASCII(b []byte) string {
	return string(b)
}

func Hex2String(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

func Hex2Str(b []byte) string {
	return Hex2String(b)
}

func HexRev2String(b []byte) string {
	return Hex2String(Rev(b))
}

func HexRev2Str(b []byte) string {
	return Hex2String(Rev(b))
}

func Uint322Hex(ui uint32) []byte {
	answer := make([]byte, 4)
	answer[3] = uint8(ui % 256)
	ui /= 256
	answer[2] = uint8(ui % 256)
	ui /= 256
	answer[1] = uint8(ui % 256)
	ui /= 256
	answer[0] = uint8(ui % 256)

	return answer
}

func Uint322HexRev(ui uint32) []byte {
	answer := make([]byte, 4)
	answer[0] = uint8(ui % 256)
	ui /= 256
	answer[1] = uint8(ui % 256)
	ui /= 256
	answer[2] = uint8(ui % 256)
	ui /= 256
	answer[3] = uint8(ui % 256)

	return answer
}

func Uint642Hex(ui uint64) []byte {
	answer := make([]byte, 8)
	answer[7] = uint8(ui % 256)
	ui /= 256
	answer[6] = uint8(ui % 256)
	ui /= 256
	answer[5] = uint8(ui % 256)
	ui /= 256
	answer[4] = uint8(ui % 256)
	ui /= 256
	answer[3] = uint8(ui % 256)
	ui /= 256
	answer[2] = uint8(ui % 256)
	ui /= 256
	answer[1] = uint8(ui % 256)
	ui /= 256
	answer[0] = uint8(ui % 256)

	return answer
}

func Uint642HexRev(ui uint64) []byte {
	answer := make([]byte, 8)
	answer[0] = uint8(ui % 256)
	ui /= 256
	answer[1] = uint8(ui % 256)
	ui /= 256
	answer[2] = uint8(ui % 256)
	ui /= 256
	answer[3] = uint8(ui % 256)
	ui /= 256
	answer[4] = uint8(ui % 256)
	ui /= 256
	answer[5] = uint8(ui % 256)
	ui /= 256
	answer[6] = uint8(ui % 256)
	ui /= 256
	answer[7] = uint8(ui % 256)

	return answer
}

func Uint162Hex(ui uint16) []byte {
	answer := make([]byte, 2)
	answer[1] = uint8(ui % 256)
	ui /= 256
	answer[0] = uint8(ui % 256)

	return answer
}

func Uint162HexRev(ui uint16) []byte {
	answer := make([]byte, 2)
	answer[0] = uint8(ui % 256)
	ui /= 256
	answer[1] = uint8(ui % 256)

	return answer
}

//TODO: test
func Uint2Hex(ui uint) []byte {
	length := 1
	if ui > 1 {
		length = int(math.Ceil(math.Log2(float64(ui)) / 8.0))
	}
	//log.Printf("%d - ui, %d - len", ui, length)
	answer := make([]byte, length)
	tmp := ui
	for i := 0; i < length; i++ {
		answer[length-1-i] = uint8(tmp % 256)
		tmp = tmp / 256
	}
	return answer
}

//TODO: test
func Int2BitHex(i int) []byte { //for the Bitcoin Script
	var ui uint
	if i < 0 {
		ui = uint(-i)
	} else {
		ui = uint(i)
	}
	answer := Uint2Hex(ui)

	if i < 0 {
		if answer[0] > 0x7F {
			answer = append([]byte{0x80}, answer[:]...)
		} else {
			answer[0] += 0x80
		}
	} else {
		if answer[0] > 0x7F {
			answer = append([]byte{0x00}, answer[:]...)
		}
	}
	return answer
}

func Int642Hex(i64 int64) []byte {
	length := 1
	if i64 > 1 {
		length = int(math.Ceil(math.Log2(float64(i64)) / 8.0))
	}
	//log.Printf("%d - ui, %d - len", ui, length)
	answer := make([]byte, length)
	tmp := i64
	for i := 0; i < length; i++ {
		answer[length-1-i] = uint8(tmp % 256)
		tmp = tmp / 256
	}
	return answer
}

func Hex2Int64(b []byte) int64 {
	var answer int64
	answer = 0
	if len(b) > 0 {
		maxBytes := len(b)
		for i := 0; i < maxBytes; i++ {
			answer *= 256
			answer += int64(b[i])
		}
	}
	return answer
}

func Byte2String(b []byte) string {
	return bytes.NewBuffer(b).String()
}

func OneByte2String(b byte) string {
	return fmt.Sprintf("%X", b)
	//return bytes.NewBuffer([]byte{b}).String()
}

func FactorToPercentageString(factor *inf.Dec) string {
	f := Str2Float(factor.String())
	f *= 100
	return Float2Str(f)
}
