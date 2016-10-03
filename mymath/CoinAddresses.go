package mymath

// Copyright 2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"crypto/rand"
	"github.com/ThePiachu/Go/Log"
	"github.com/ThePiachu/Go/mymath/bitecdsa"
	"github.com/ThePiachu/Go/mymath/bitelliptic"
	"golang.org/x/net/context"
)

type CoinAddress struct {
	NetByte               string
	PrivateKey            string
	PublicKeyCompressed   string
	PublicKeyUncompressed string
	AddressCompressed     string
	AddressUncompressed   string
}

func (ca *CoinAddress) GetWIFCompressed(prefixByte string) string {
	privateKey := ca.PrivateKey + "01"
	pb := String2Hex(prefixByte)[0]
	wif := PrivateKeyToWIFWithPrefixByte(privateKey, pb)
	return wif
}
func (ca *CoinAddress) GetWIFUncompressed(prefixByte string) string {
	privateKey := ca.PrivateKey
	pb := String2Hex(prefixByte)[0]
	wif := PrivateKeyToWIFWithPrefixByte(privateKey, pb)
	return wif
}

func GenerateNewCoinAddress(netByte string) CoinAddress {
	address := CoinAddress{}
	address.NetByte = netByte
	curve := bitelliptic.S256()
	priv, err := bitecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return address
	}
	address.PrivateKey = Hex2Str(Big2Hex(priv.D))

	address.PublicKeyCompressed = PrivateKeyToCompressedPublicKey(address.PrivateKey)
	address.PublicKeyUncompressed = PrivateKeyToUncompressedPublicKey(address.PrivateKey)
	address.AddressCompressed = PublicKeyToAddress(netByte, address.PublicKeyCompressed)
	address.AddressUncompressed = PublicKeyToAddress(netByte, address.PublicKeyUncompressed)
	return address
}

func GenerateUncompressedCoinVanityAddressWithNetbyte(c context.Context, pattern string, netbyte string) CoinAddress {
	subpattern := pattern[0 : len(pattern)-1]
	Log.Infof(c, "Subpattern - %v", subpattern)
	Log.Infof(c, "Pattern - %v", pattern)
	for i := 0; ; i++ {
		address := GenerateNewCoinAddress(netbyte)
		if DoesStringStartWith(address.AddressUncompressed, subpattern) {
			Log.Infof(c, "Trying %v - %v", i, address.AddressUncompressed)
			Log.Infof(c, "%v", address)

			if DoesStringStartWith(address.AddressUncompressed, pattern) {
				return address
			}
		}
		if i%1000 == 0 {
			Log.Infof(c, "Trying %v - %v", i, address.AddressUncompressed)
		}
	}
	return CoinAddress{}
}

func GenerateCompressedCoinVanityAddressWithNetbyte(c context.Context, pattern string, netbyte string) CoinAddress {
	subpattern := pattern[0 : len(pattern)-1]
	Log.Infof(c, "Subpattern - %v", subpattern)
	Log.Infof(c, "Pattern - %v", pattern)
	for i := 0; ; i++ {
		address := GenerateNewCoinAddress(netbyte)
		if DoesStringStartWith(address.AddressCompressed, subpattern) {
			Log.Infof(c, "Trying %v - %v", i, address.AddressCompressed)
			Log.Infof(c, "%v", address)

			if DoesStringStartWith(address.AddressCompressed, pattern) {
				return address
			}
		}
		if i%1000 == 0 {
			Log.Infof(c, "Trying %v - %v", i, address.AddressCompressed)
		}
	}
	return CoinAddress{}
}

func GenerateCoinVanityAddressWithNetbyte(c context.Context, pattern string, netbyte string) CoinAddress {
	subpattern := pattern[0 : len(pattern)-1]
	Log.Infof(c, "Subpattern - %v", subpattern)
	Log.Infof(c, "Pattern - %v", pattern)
	for i := 0; ; i++ {
		address := GenerateNewCoinAddress(netbyte)
		if DoesStringStartWith(address.AddressCompressed, subpattern) {
			Log.Infof(c, "Trying %v - %v", i, address.AddressCompressed)
			Log.Infof(c, "%v", address)

			if DoesStringStartWith(address.AddressCompressed, pattern) {
				return address
			}
		}
		if DoesStringStartWith(address.AddressUncompressed, subpattern) {
			Log.Infof(c, "Trying %v - %v", i, address.AddressUncompressed)
			Log.Infof(c, "%v", address)

			if DoesStringStartWith(address.AddressUncompressed, pattern) {
				return address
			}
		}
		if i%1000 == 0 {
			Log.Infof(c, "Trying %v - %v", i, address.AddressUncompressed)
		}
	}
	return CoinAddress{}
}
