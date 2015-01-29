package mymath

// Copyright 2011-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package for handling common math and conversion operations used in the rest of the program.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
	"strings"
)

func IsStringAValidEmailAddress(s string) bool {
	//TODO: expand
	if strings.Count(s, "@") != 1 {
		return false
	}
	if strings.ContainsAny(s, " \"(),:;<>[]\\") {
		return false
	}

	return true
}

func DecodeJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, &v)
	return err
}

func EncodeJSON(data interface{}) ([]byte, error) {
	encoded, err := json.Marshal(data)
	return encoded, err
}

func EncodeJSONString(data interface{}) (string, error) {
	encoded, err := EncodeJSON(data)
	if err != nil {
		return "", err
	}
	return Hex2ASCII(encoded), err
}

func DecodeJSONString(data string, v interface{}) error {
	return DecodeJSON(ASCII2Hex(data), v)
}

func Base64URLEncode(s string) string {
	ans := base64.URLEncoding.EncodeToString([]byte(s))
	return ans
}

func Base64URLDecode(s string) string {
	ans, _ := base64.URLEncoding.DecodeString(s)
	return string(ans)
}

func SplitStrings(str, separator string) []string {
	return strings.Split(str, separator)
}

func Contains(s string, sub string) bool {
	return strings.Contains(s, sub)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func AreStringsEqual(one, two string) bool {
	if len(one) != len(two) {
		return false
	}
	for i := 0; i < len(one); i++ {
		if one[i] != two[i] {
			return false
		}
	}
	return true
}

func DoesStringStartWith(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func DoesStringContain(s string, substring string) bool {
	return strings.Contains(s, substring)
}

func AreHexesEqual(one, two []byte) bool {
	if len(one) != len(two) {
		return false
	}
	for i := 0; i < len(one); i++ {
		if one[i] != two[i] {
			return false
		}
	}
	return true
}

//TODO: use this function to reverse some of the below functions?
func Rev(b []byte) []byte {
	answer := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		answer[i] = b[len(b)-1-i]
	}
	return answer
}

func RevWords(b []byte) []byte {
	answer := make([]byte, len(b))
	for i := 0; i < len(b)/4; i++ {
		answer[i*4+3] = b[len(b)-4-i*4+3]
		answer[i*4+2] = b[len(b)-4-i*4+2]
		answer[i*4+1] = b[len(b)-4-i*4+1]
		answer[i*4+0] = b[len(b)-4-i*4+0]
	}
	return answer
}

func RevWordsStr(b string) string {
	return Hex2Str(RevWords(Str2Hex(b)))
}
func RevWords2Str(b string) string {
	return Hex2Str(RevWords2(Str2Hex(b)))
}

func RevWords2(b []byte) []byte {
	answer := make([]byte, len(b))
	for i := 0; i < len(b)/4; i++ {
		answer[i*4+3] = b[i*4+0]
		answer[i*4+2] = b[i*4+1]
		answer[i*4+1] = b[i*4+2]
		answer[i*4+0] = b[i*4+3]
	}
	return answer
}

func Hex2Big(b []byte) *big.Int {
	answer := big.NewInt(0)

	for i := 0; i < len(b); i++ {
		answer.Lsh(answer, 8)
		answer.Add(answer, big.NewInt(int64(b[i])))
	}

	return answer
}

func HexRev2Big(rev []byte) *big.Int {

	b := make([]byte, len(rev))

	for i := 0; i < len(rev); i++ {
		b[len(rev)-i-1] = rev[i]
	}

	answer := big.NewInt(0)

	for i := 0; i < len(b); i++ {
		answer.Lsh(answer, 8)
		answer.Add(answer, big.NewInt(int64(b[i])))
	}

	return answer
}

func Big2HexPadded(b *big.Int, paddingSize int) []byte {
	answer := Big2Hex(b)
	//log.Printf("len answer - %d", len(answer))
	if len(answer) == paddingSize {
		return answer
	}
	if len(answer) > paddingSize {
		return answer
	}
	return append(make([]byte, paddingSize-len(answer)), answer...)
}

func ConcatBytes(list ...[]byte) []byte {
	size := 0 //size of the resulting concatenated list
	for i := 0; i < len(list); i++ {
		size += len(list[i]) //counting the sizes of individual parts of the list
	}
	answer := make([]byte, size) //creates the array for the answer

	iterator := 0 //iterator to count the position in the answer array
	for i := 0; i < len(list); i++ {
		copy(answer[iterator:], list[i]) //copies the data into the answer array
		iterator += len(list[i])
	}
	return answer //returns the result
}

func AddByte(one []byte, two byte) []byte {
	size := len(one) + 1 //size of the resulting concatenated list

	answer := make([]byte, size) //creates the array for the answer

	copy(answer[0:], one) //copies the data into the answer array
	answer[len(one)] = two
	return answer //returns the result
}

/*func SqrtBigInt(n *big.Int) *big.Int {
	//http://www.codecodex.com/wiki/Calculate_an_integer_square_root
	one:=big.NewInt(1)
	two:=big.NewInt(2)
	xn:=big.NewInt(1)
	xn1:=big.NewInt(1)
	xn1.Add(xn1, n)
	xn1.Div(xn1, two)
	diff:=big.NewInt(0)
	diff.Sub(xn1, xn)
	diff.Abs(diff)
	for diff.Cmp(one)>0 {
		xn.Set(xn1)
		xn1.Div(n, xn)
		xn1.Add(xn, xn1)
		xn1.Div(xn1, two)

		diff.Sub(xn1, xn)
		diff.Abs(diff)
	}
	xn2:=big.NewInt(1)
	xn2.Exp(xn1, two, nil)
	for xn2.Cmp(n)>0 {
		xn1.Sub(xn1, one)
		xn2.Exp(xn1, two, nil)
	}
	return xn2
}*/

//Testing

//TODO: do
func TestEverything() bool {
	//TestEverythingBitmath()

	if RevTest() == false {
		return false
	}

	log.Print("All tests okay!")
	return true
}

func RevTest() bool {
	one := make([]byte, 3)
	two := make([]byte, 3)
	one[0] = 0xFE
	one[1] = 0xA9
	one[2] = 0x01

	two[0] = 0x01
	two[1] = 0xA9
	two[2] = 0xFE
	if bytes.Compare(Rev(one), two) != 0 {
		return false
	}
	return true
}
