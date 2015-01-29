package mymath

// Copyright 2013-2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"appengine"
	"github.com/ThePiachu/Go/Log"
)

var RippleAlphabet = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"

type RippleAddress struct {
	AccountID     string
	MasterSeed    string
	MasterSeedHex string

	PrivateGenerator   string
	PublicGenerator    string
	PublicGeneratorHex string

	PrivateKey   string
	PublicKey    string
	PublicKeyHex string
}

var RippleRevAlp = map[string]int{
	"r": 0, "p": 1, "s": 2, "h": 3, "n": 4, "a": 5, "f": 6, "3": 7, "9": 8, "w": 9,
	"B": 10, "U": 11, "D": 12, "N": 13, "E": 14, "G": 15, "H": 16, "J": 17, "K": 18, "L": 19,
	"M": 20, "4": 21, "P": 22, "Q": 23, "R": 24, "S": 25, "T": 26, "7": 27, "V": 28, "W": 29,
	"X": 30, "Y": 31, "Z": 32, "2": 33, "b": 34, "c": 35, "d": 36, "e": 37, "C": 38, "g": 39,
	"6": 40, "5": 41, "j": 42, "k": 43, "m": 44, "8": 45, "o": 46, "F": 47, "q": 48, "i": 49,
	"1": 50, "t": 51, "u": 52, "v": 53, "A": 54, "x": 55, "y": 56, "z": 57,
}

func RippleAddressToStandardAddress(address string) string {
	answer := ""
	for i := 0; i < len(address); i++ {
		letter := address[i : i+1]
		index, ok := RippleRevAlp[letter]
		if ok == false {
			return ""
		}
		answer += BitcoinAlphabet[index : index+1]
	}
	return answer
}

func StandardAddressToRippleAddress(address string) string {
	answer := ""
	for i := 0; i < len(address); i++ {
		letter := address[i : i+1]
		index, ok := BitcoinRevAlp[letter]
		if ok == false {
			return ""
		}
		answer += RippleAlphabet[index : index+1]
	}
	return answer
}

func CheckRippleAddressValidity(address string) bool {
	standardAddress := RippleAddressToStandardAddress(address)
	return CheckAddressStringValidity(standardAddress)
}

func CheckRippleSecretValidity(secret string, address string) bool {
	//TODO: do

	return true
}

func GenerateBase58CheckString(baseHex string, leadingByte string) string {
	base := append(Str2Hex(leadingByte), Str2Hex(baseHex)...)
	tmp := DoubleSHA(base)
	base = append(base, tmp[0:4]...)
	return Hex2Str(base)
}

func GenerateNewRippleAddress() RippleAddress {
	address := RippleAddress{}
	//address.MasterSeedHex="71ED064155FFADFA38782C5E0158CB26"
	address.MasterSeedHex = Hex2Str(RandomHex(16))
	base := string(StrHex2Base58(GenerateBase58CheckString(address.MasterSeedHex, "21")))
	address.MasterSeed = StandardAddressToRippleAddress(base)

	privateGenerator := ""
	var sequence uint32
	for i := 0; ; i++ {
		sequence = uint32(i)
		temp := address.MasterSeedHex + Hex2Str(Uint322Hex(sequence))
		privateGenerator = SHA512HalfString(temp)
		if IsPrivateKeyOnCurve(privateGenerator) == true {
			break
		}
	}
	address.PrivateGenerator = privateGenerator
	//log.Printf("privateGenerator - %v", privateGenerator)
	address.PublicGeneratorHex = PrivateKeyToCompressedPublicKey(address.PrivateGenerator)
	address.PublicGenerator = StandardAddressToRippleAddress(string(StrHex2Base58(GenerateBase58CheckString(address.PublicGeneratorHex, "29"))))
	//log.Printf("address.PublicGeneratorHex - %v", address.PublicGeneratorHex)

	var subsequence uint32
	privateGenerator2 := ""
	for i := 0; ; i++ {
		subsequence = uint32(i)
		temp := address.PublicGeneratorHex + Hex2Str(Uint322Hex(sequence)) + Hex2Str(Uint322Hex(subsequence))
		//log.Printf("temp - %v", temp)
		privateGenerator2 = SHA512HalfString(temp)
		//log.Printf("privateGenerator2 - %v", privateGenerator2)
		if IsPrivateKeyOnCurve(privateGenerator2) == true {
			//if ComparePrivateKeys(privateGenerator, privateGenerator2)>=0 {
			break
			//}
		}
	}
	combined := AddPrivateKeysReturnString(privateGenerator, privateGenerator2)
	//log.Printf("combined - %v", combined)
	address.PrivateKey = combined
	address.PublicKeyHex = PrivateKeyToCompressedPublicKey(combined)
	address.PublicKey = StandardAddressToRippleAddress(string(StrHex2Base58(GenerateBase58CheckString(address.PublicKeyHex, "23"))))

	rippleAddress := PublicKeyToAddress("00", address.PublicKeyHex)

	address.AccountID = StandardAddressToRippleAddress(rippleAddress)

	return address
}

func GenerateRippleVanityAddress(c appengine.Context, pattern string) RippleAddress {
	subpattern := pattern[0 : len(pattern)-1]
	Log.Infof(c, "Subpattern - %v", subpattern)
	Log.Infof(c, "Pattern - %v", pattern)
	for i := 0; ; i++ {
		address := GenerateNewRippleAddress()
		if DoesStringStartWith(address.AccountID, subpattern) {
			Log.Infof(c, "Trying %v - %v", i, address.AccountID)
			Log.Infof(c, "%v", address)

			if DoesStringStartWith(address.AccountID, pattern) {
				return address
			}
		}
		if i%1000 == 0 {
			Log.Infof(c, "Trying %v - %v", i, address.AccountID)
		}
	}
	return RippleAddress{}
}

func ComparePrivateKeys(a string, b string) int {
	aBig := Str2Big(a)
	bBig := Str2Big(b)
	return aBig.Cmp(bBig)
}
