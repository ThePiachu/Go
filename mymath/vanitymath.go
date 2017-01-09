package mymath

// Copyright 2012 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ThePiachu/Go/mymath/bitecdsa"
	"github.com/ThePiachu/Go/mymath/bitelliptic"
)

func IsPublicKeyValid(pubKey string) bool {
	if pubKey == "" {
		return false
	}
	if len(pubKey) != 130 {
		return false
	}
	if pubKey[0] != '0' || pubKey[1] != '4' {
		return false
	}

	a, b := PublicKeyToPointCoordinates(pubKey)
	if a == nil || b == nil {
		return false
	}

	curve := bitelliptic.S256()
	if !curve.IsOnCurve(a, b) {
		return false
	}

	return true
}

func VanityAddressLavishness(pattern string, bounty float64) float64 {
	complexity := VanityAddressComplexity(pattern)
	return CalculateLavishness(bounty, complexity)
}

func CalculateBountyForLavishnessAndComplexity(lavishness float64, complexity float64) float64 {
	return lavishness * complexity / math.Exp2(32.0)
}

func VanityAddressComplexity(pattern string) float64 {
	if len(pattern) < 1 {
		return 0
	}
	complexity := 1.0

	countingOnes := true
	if len(pattern) == 1 {
		return 1.0
	}
	switch pattern[1] {
	case '1':
		complexity *= 256
	case '2', '3', '4', '5', '6':
		complexity *= 23
		countingOnes = false
	case '7', '8', '9', 'A',
		'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L',
		'M', 'N', 'P':
		complexity *= 22
		countingOnes = false
	case 'Q':
		complexity *= 65
		countingOnes = false
	case 'R', 'S', 'T', 'U', 'V', 'W',
		'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g',
		'h', 'i', 'j', 'k', 'm', 'n', 'o', 'p', 'q', 'r',
		's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		complexity *= 1353
		countingOnes = false
	}

	for i := 2; i < len(pattern); i++ {
		if countingOnes {
			if pattern[i] == '1' {
				complexity *= 256
			} else {
				complexity *= 58
				countingOnes = false
			}
		} else {
			complexity *= 58
		}
	}

	return complexity
}

func CalculateLavishness(bounty float64, complexity float64) float64 {
	return math.Exp2(32.0) * bounty / complexity
}

func AddPrivateKeys(one, two string) *big.Int {
	a := big.NewInt(0)
	b := big.NewInt(0)

	a.SetString(one, 16)
	b.SetString(two, 16)

	one1 := big.NewInt(1)

	tmp := a.Add(a, b)
	tmp = tmp.Sub(tmp, one1)
	//mod,_:=new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
	//mod:=new(big.Int).Sub(bitelliptic.S256().P, one1)
	mod := bitelliptic.S256().N
	tmp = tmp.Mod(tmp, mod)
	answer := tmp.Add(tmp, one1)
	//answer:=tmp
	//answer:=a.Mod(a.Add(a, b), bitelliptic.S256().P)

	//log.Printf("Sum 1 - %X", Big2Hex(answer))

	return answer
}

func AddPrivateKeysReturnString(one, two string) string {
	sum := AddPrivateKeys(one, two)
	return Hex2Str(Big2HexWithMinimumLength(sum, 32))
}

func MultiplyPrivateKeys(one, two string) *big.Int {
	a := big.NewInt(0)
	b := big.NewInt(0)

	a.SetString(one, 16)
	b.SetString(two, 16)

	one1 := big.NewInt(1)

	tmp := a.Mul(a, b)
	tmp = tmp.Sub(tmp, one1)
	//mod,_:=new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
	//mod:=new(big.Int).Sub(bitelliptic.S256().P, one1)
	mod := bitelliptic.S256().N
	tmp = tmp.Mod(tmp, mod)
	answer := tmp.Add(tmp, one1)
	//answer:=tmp
	//answer:=a.Mod(a.Add(a, b), bitelliptic.S256().P)

	return answer
}
func MultiplyPrivateKeysReturnString(one, two string) string {
	mult := MultiplyPrivateKeys(one, two)
	return Hex2Str(Big2HexWithMinimumLength(mult, 32))
}

func AddPublicKeys(one, two string) (*big.Int, *big.Int) {
	a, b := PublicKeyToPointCoordinates(one)
	c, d := PublicKeyToPointCoordinates(two)
	e, f := bitelliptic.S256().Add(a, b, c, d)

	log.Printf("04%X%X", Big2Hex(e), Big2Hex(f))

	return e, f
}

func AddPublicKeysReturnString(one, two string) string {
	a, b := AddPublicKeys(one, two)
	return PointCoordinatesToPublicKey(a, b)
}

func MultiplyPrivateAndPublicKeyReturnString(private, public string) string {
	gx, gy := PublicKeyToPointCoordinates(public)

	curve2 := bitelliptic.S256().Copy()
	curve2.Gx = gx
	curve2.Gy = gy

	priv2, _ := bitecdsa.GenerateFromPrivateKey(Str2Big(private), curve2)

	//priv, _:=bitecdsa.GenerateFromPrivateKey(priv2.D, curve2)

	answer := append(append([]byte{0x04}, Big2Hex(priv2.PublicKey.X)...), Big2Hex(priv2.PublicKey.Y)...)

	return Hex2Str(answer)
}

func PrivateKeyToWIF(key string) string {
	priv := append([]byte{0x80}, Str2Hex(key)...)
	sha := SingleSHA(priv)
	sha = SingleSHA(sha)
	priv = append(priv, sha[0:4]...)
	return string(Hex2Base58(priv))
}

func PrivateKeyToWIFWithPrefixByte(key string, prefix byte) string {
	priv := append([]byte{prefix}, Str2Hex(key)...)
	sha := SingleSHA(priv)
	sha = SingleSHA(sha)
	priv = append(priv, sha[0:4]...)
	return string(Hex2Base58(priv))
}

func PublicKeyToPointCoordinates(pubKey string) (*big.Int, *big.Int) {
	if len(pubKey) != 130 {
		log.Printf("PublicKeyToPointCoordinates - len(pubKey)!=130")
		log.Printf(pubKey)
		return nil, nil
	}
	if pubKey[0] != '0' || pubKey[1] != '4' {
		log.Printf("pubKey[0]!='0' || pubKey[1]!='4'")
		return nil, nil
	}
	a := big.NewInt(0)
	b := big.NewInt(0)

	a.SetString(pubKey[2:66], 16)
	b.SetString(pubKey[66:130], 16)

	log.Printf("pubKey - %s", pubKey)
	log.Printf("X - %v, Y - %s", a.String(), b.String())

	return a, b
}

func PointCoordinatesToPublicKey(x, y *big.Int) string {
	one := Big2Hex(x)
	if len(one) < 32 {
		tmp := make([]byte, 32-len(one))
		one = append(tmp, one...)
	}
	two := Big2Hex(y)
	if len(two) < 32 {
		tmp := make([]byte, 32-len(two))
		two = append(tmp, two...)
	}
	return "04" + Hex2Str(one) + Hex2Str(two)
}

func CheckSolutionAddUncompressed(pubKey, solution, pattern string, netByte byte) (string, string) {
	if IsPrivateKeyOnCurve(solution) == false {
		return "", "Point is not on curve"
	}
	pubKey2 := PrivateKeyToUncompressedPublicKey(solution)
	sum := AddPublicKeysReturnString(pubKey, pubKey2)
	address := PublicKeyToAddress(OneByte2String(netByte), sum)
	for i := 0; i < len(pattern); i++ {
		if address[i] != pattern[i] {
			return "", "Wrong pattern"
		}
	}
	return address, ""
}

func CheckSolutionAddCompressed(pubKey, solution, pattern string, netByte byte) (string, string) {
	if IsPrivateKeyOnCurve(solution) == false {
		return "", "Point is not on curve"
	}
	pubKey2 := PrivateKeyToUncompressedPublicKey(solution)
	sum := UncompressedKeyToCompressedKey(AddPublicKeysReturnString(pubKey, pubKey2))
	address := PublicKeyToAddress(OneByte2String(netByte), sum)
	for i := 0; i < len(pattern); i++ {
		if address[i] != pattern[i] {
			return "", "Wrong pattern"
		}
	}
	return address, ""
}

/*
func CheckSolutionAdd(pubKey, solution, pattern string, netByte byte) (string, string){
	d:=big.NewInt(0)
	d.SetString(solution, 16)

	private, err:=bitecdsa.GenerateFromPrivateKey(d, bitelliptic.S256())
	if err!=nil{
		log.Printf("vanitymath err - %s", err)
		return "", err.Error()
	}
	a, b:=PublicKeyToPointCoordinates(pubKey)

	x, y:=bitelliptic.S256().Add(a, b, private.PublicKey.X, private.PublicKey.Y)

	ba:=NewFromPublicKeyString(netByte, PointCoordinatesToPublicKey(x, y))
	address:=string(ba.Base)
	for i:=0;i<len(pattern);i++{
		if address[i]!=pattern[i]{
			return "", "Wrong pattern"
		}
	}
	return address, ""
}*/

func CheckSolutionMultUncompressed(pubKey, solution, pattern string, netByte byte) (string, string) {
	if IsPrivateKeyOnCurve(solution) == false {
		return "", "Point is not on curve"
	}
	mult := MultiplyPrivateAndPublicKeyReturnString(solution, pubKey)
	address := PublicKeyToAddress(OneByte2String(netByte), mult)
	for i := 0; i < len(pattern); i++ {
		if address[i] != pattern[i] {
			return "", "Wrong pattern"
		}
	}
	return address, ""
}

func CheckSolutionMultCompressed(pubKey, solution, pattern string, netByte byte) (string, string) {
	if IsPrivateKeyOnCurve(solution) == false {
		return "", "Point is not on curve"
	}
	mult := MultiplyPrivateAndPublicKeyReturnString(solution, pubKey)
	address := PublicKeyToAddress(OneByte2String(netByte), UncompressedKeyToCompressedKey(mult))
	for i := 0; i < len(pattern); i++ {
		if address[i] != pattern[i] {
			return "", "Wrong pattern"
		}
	}
	return address, ""
}

/*func CheckSolutionMult(pubKey, solution, pattern string, netByte byte) (string, string){
	gx, gy:=PublicKeyToPointCoordinates(pubKey)

	curve2:=bitelliptic.S256().Copy()
	curve2.Gx=gx
	curve2.Gy=gy

	ba, err:=NewAddressFromPrivateKeyWithCurveOtherNets(netByte, Str2Hex(solution), curve2)

	if err!=nil{
		return "", err.Error()
	}

	address:=string(ba.Base)
	for i:=0;i<len(pattern);i++{
		if address[i]!=pattern[i]{
			return "", "Wrong pattern"
		}
	}
	return address, ""
}*/

func NewRandomPrivateKey() string {
	curve := bitelliptic.S256()
	priv, _, _, _ := curve.GenerateKey(rand.Reader)
	return Hex2Str(priv)
}

func PrivateKeyToCompressedAddressBytes(net byte, privateKey []byte) string {
	return PrivateKeyToCompressedAddress(fmt.Sprintf("%x", net), fmt.Sprintf("%x", privateKey))
}

func PrivateKeyToUncompressedAddressBytes(net byte, privateKey []byte) string {
	return PrivateKeyToUncompressedAddress(fmt.Sprintf("%x", net), fmt.Sprintf("%x", privateKey))
}

func PrivateKeyToCompressedAddress(net string, privateKey string) string {
	pub := PrivateKeyToCompressedPublicKey(privateKey)
	return PublicKeyToAddress(net, pub)
}

func PrivateKeyToUncompressedAddress(net string, privateKey string) string {
	pub := PrivateKeyToUncompressedPublicKey(privateKey)
	return PublicKeyToAddress(net, pub)
}

func PrivateKeyToCompressedPublicKey(privateKey string) string {
	uncompressed := PrivateKeyToUncompressedPublicKey(privateKey)
	return UncompressedKeyToCompressedKey(uncompressed)
}

func PrivateKeyToUncompressedPublicKey(privateKey string) string {
	curve := bitelliptic.S256()
	priv, err := bitecdsa.GenerateFromPrivateKey(Str2Big(privateKey), curve)
	if err != nil {
		return ""
	}
	x := Big2Hex(priv.PublicKey.X)
	for len(x) < 32 {
		x = append([]byte{0x00}, x...)
	}
	y := Big2Hex(priv.PublicKey.Y)
	for len(y) < 32 {
		y = append([]byte{0x00}, y...)
	}
	publicKey := append(append([]byte{0x04}, x...), y...)
	return Hex2Str(publicKey)
}

func UncompressedKeyToCompressedKey(key string) string {
	if len(key) != 130 {
		return ""
	}
	x := key[2:66]
	y := HexString2Int64(key[129:])
	if y%2 == 0 {
		return "02" + x
	} else {
		return "03" + x
	}
}

func PublicKeyToAddress(net string, key string) string {
	if len(net) == 0 {
		return ""
	}
	netByte := Str2Hex(net)
	keyBytes := Str2Hex(key)

	hash := make([]byte, 25)
	hash[0] = netByte[0]

	tmp := SHARipemd(keyBytes)

	for i := 0; i < len(tmp); i++ {
		hash[i+1] = tmp[i]
	}

	tmp = DoubleSHA(hash[0:21])

	//hash checksum
	hash[21] = tmp[0]
	hash[22] = tmp[1]
	hash[23] = tmp[2]
	hash[24] = tmp[3]

	return string(Hex2Base58(hash))
}

/*func CompressedKeyToUncompressedKey(key string) string {
	//Y^2=X^3+aX+b

	//a=0, defined by curve
	b:=big.NewInt(7)//defined by curve

	//Y^2=X^3+b
	//Y=+-sqrt(X^3+b)

	x:=String2Big(key[2:66])
	x.Exp(x, big.NewInt(3), nil)
	x.Add(x, b)

	y:=SqrtBigInt(x)
	curve:=bitelliptic.S256()
	y.Neg(y)
	y.Mod(y, curve.P)

	return "04"+key[2:66]+Hex2Str(Big2HexWithMinimumLength(y, 32))
}*/

func IsPrivateKeyOnCurve(privateKey string) bool {
	curve := bitelliptic.S256()
	return bitecdsa.CheckIsOnCurve(curve, Str2Big(privateKey))
}

func DoesPatternHaveRightPrefix(pattern string, netByte byte) bool {
	if len(pattern) == 0 {
		return false
	}
	if WillNetworkByteGivePrefix(netByte, pattern[0:1]) == false {
		return false
	}
	return CanPatternBeSolvedInNet(pattern, netByte)
}

func MinAddressForNetByte(netByte byte) string {
	minBytes := make([]byte, 25)
	minBytes[0] = netByte
	for i := 1; i < 25; i++ {
		minBytes[i] = 0x00
	}
	return string(Hex2Base58(minBytes))
}

func MaxAddressForNetByte(netByte byte) string {
	maxBytes := make([]byte, 25)
	maxBytes[0] = netByte
	for i := 1; i < 25; i++ {
		maxBytes[i] = 0xFF
	}
	return string(Hex2Base58(maxBytes))
}

func CanPatternBeSolvedInNet(pattern string, netByte byte) bool {
	minAddress := MinAddressForNetByte(netByte)
	maxAddress := MaxAddressForNetByte(netByte)

	log.Printf("pattern - %v, minAddress - %v, maxAddress - %v", pattern, minAddress, maxAddress)

	if netByte == 0 && pattern[0] == '1' {
		//mainnet is harder to check but any pattern is solvable
		return true
	}

	for i := 0; i < len(pattern); i++ {
		p := Base58LetterToInt(pattern[i : i+1])
		min := Base58LetterToInt(minAddress[i : i+1])
		log.Printf("p - %v, min - %v", p, min)
		if p < min {
			return false
		}
		if p > min {
			break
		}
	}
	for i := 0; i < len(pattern); i++ {
		p := Base58LetterToInt(pattern[i : i+1])
		max := Base58LetterToInt(maxAddress[i : i+1])
		log.Printf("p - %v, max - %v", p, max)
		if p > max {
			return false
		}
		if p < max {
			break
		}
	}
	return true
}

func WillNetworkByteGivePrefix(b byte, prefix string) bool {
	log.Printf("WillNetworkByteGivePrefix - %v, %v", b, prefix)

	switch b {
	case 0:
		if prefix == "1" {
			return true
		}
	case 1:
		if prefix == "Q" {
			return true
		}
		if prefix == "R" {
			return true
		}
		if prefix == "S" {
			return true
		}
		if prefix == "T" {
			return true
		}
		if prefix == "U" {
			return true
		}
		if prefix == "V" {
			return true
		}
		if prefix == "W" {
			return true
		}
		if prefix == "X" {
			return true
		}
		if prefix == "Y" {
			return true
		}
		if prefix == "Z" {
			return true
		}
		if prefix == "a" {
			return true
		}
		if prefix == "b" {
			return true
		}
		if prefix == "c" {
			return true
		}
		if prefix == "d" {
			return true
		}
		if prefix == "e" {
			return true
		}
		if prefix == "f" {
			return true
		}
		if prefix == "g" {
			return true
		}
		if prefix == "h" {
			return true
		}
		if prefix == "i" {
			return true
		}
		if prefix == "j" {
			return true
		}
		if prefix == "k" {
			return true
		}
		if prefix == "m" {
			return true
		}
		if prefix == "n" {
			return true
		}
		if prefix == "o" {
			return true
		}
	case 2:
		if prefix == "o" {
			return true
		}
		if prefix == "p" {
			return true
		}
		if prefix == "q" {
			return true
		}
		if prefix == "r" {
			return true
		}
		if prefix == "s" {
			return true
		}
		if prefix == "t" {
			return true
		}
		if prefix == "u" {
			return true
		}
		if prefix == "v" {
			return true
		}
		if prefix == "w" {
			return true
		}
		if prefix == "x" {
			return true
		}
		if prefix == "y" {
			return true
		}
		if prefix == "z" {
			return true
		}
	case 3:
		if prefix == "2" {
			return true
		}
	case 4:
		if prefix == "2" {
			return true
		}
		if prefix == "3" {
			return true
		}
	case 5, 6:
		if prefix == "3" {
			return true
		}
	case 7:
		if prefix == "4" {
			return true
		}
		if prefix == "3" {
			return true
		}
	case 8:
		if prefix == "4" {
			return true
		}
	case 9:
		if prefix == "4" {
			return true
		}
		if prefix == "5" {
			return true
		}
	case 10, 11:
		if prefix == "5" {
			return true
		}
	case 12:
		if prefix == "5" {
			return true
		}
		if prefix == "6" {
			return true
		}
	case 13:
		if prefix == "6" {
			return true
		}
	case 14:
		if prefix == "6" {
			return true
		}
		if prefix == "7" {
			return true
		}
	case 15, 16:
		if prefix == "7" {
			return true
		}
	case 17:
		if prefix == "7" {
			return true
		}
		if prefix == "8" {
			return true
		}
	case 18:
		if prefix == "8" {
			return true
		}
	case 19:
		if prefix == "8" {
			return true
		}
		if prefix == "9" {
			return true
		}
	case 20, 21:
		if prefix == "9" {
			return true
		}
	case 22:
		if prefix == "9" {
			return true
		}
		if prefix == "A" {
			return true
		}
	case 23:
		if prefix == "A" {
			return true
		}
	case 24:
		if prefix == "A" {
			return true
		}
		if prefix == "B" {
			return true
		}
	case 25, 26:
		if prefix == "B" {
			return true
		}
	case 27:
		if prefix == "B" {
			return true
		}
		if prefix == "C" {
			return true
		}
	case 28:
		if prefix == "C" {
			return true
		}
	case 29:
		if prefix == "C" {
			return true
		}
		if prefix == "D" {
			return true
		}
	case 30, 31:
		if prefix == "D" {
			return true
		}
	case 32:
		if prefix == "D" {
			return true
		}
		if prefix == "E" {
			return true
		}
	case 33:
		if prefix == "E" {
			return true
		}
	case 34:
		if prefix == "E" {
			return true
		}
		if prefix == "F" {
			return true
		}
	case 35, 36:
		if prefix == "F" {
			return true
		}
	case 37:
		if prefix == "F" {
			return true
		}
		if prefix == "G" {
			return true
		}
	case 38:
		if prefix == "G" {
			return true
		}
	case 39:
		if prefix == "G" {
			return true
		}
		if prefix == "H" {
			return true
		}
	case 40, 41:
		if prefix == "H" {
			return true
		}
	case 42:
		if prefix == "H" {
			return true
		}
		if prefix == "J" {
			return true
		}
	case 43:
		if prefix == "J" {
			return true
		}
	case 44:
		if prefix == "J" {
			return true
		}
		if prefix == "K" {
			return true
		}
	case 45, 46:
		if prefix == "K" {
			return true
		}
	case 47:
		if prefix == "K" {
			return true
		}
		if prefix == "L" {
			return true
		}
	case 48:
		if prefix == "L" {
			return true
		}
	case 49:
		if prefix == "L" {
			return true
		}
		if prefix == "M" {
			return true
		}
	case 50, 51:
		if prefix == "M" {
			return true
		}
	case 52:
		if prefix == "M" {
			return true
		}
		if prefix == "N" {
			return true
		}
	case 53:
		if prefix == "N" {
			return true
		}
	case 54:
		if prefix == "N" {
			return true
		}
		if prefix == "P" {
			return true
		}
	case 55, 56:
		if prefix == "P" {
			return true
		}
	case 57:
		if prefix == "P" {
			return true
		}
		if prefix == "Q" {
			return true
		}
	case 58:
		if prefix == "Q" {
			return true
		}
	case 59:
		if prefix == "Q" {
			return true
		}
		if prefix == "R" {
			return true
		}
	case 60, 61:
		if prefix == "R" {
			return true
		}
	case 62:
		if prefix == "R" {
			return true
		}
		if prefix == "S" {
			return true
		}
	case 63:
		if prefix == "S" {
			return true
		}
	case 64:
		if prefix == "S" {
			return true
		}
		if prefix == "T" {
			return true
		}
	case 65, 66:
		if prefix == "T" {
			return true
		}
	case 67:
		if prefix == "T" {
			return true
		}
		if prefix == "U" {
			return true
		}
	case 68:
		if prefix == "U" {
			return true
		}
	case 69:
		if prefix == "U" {
			return true
		}
		if prefix == "V" {
			return true
		}
	case 70, 71:
		if prefix == "V" {
			return true
		}
	case 72:
		if prefix == "V" {
			return true
		}
		if prefix == "W" {
			return true
		}
	case 73:
		if prefix == "W" {
			return true
		}
	case 74:
		if prefix == "W" {
			return true
		}
		if prefix == "X" {
			return true
		}
	case 75, 76:
		if prefix == "X" {
			return true
		}
	case 77:
		if prefix == "X" {
			return true
		}
		if prefix == "Y" {
			return true
		}
	case 78:
		if prefix == "Y" {
			return true
		}
	case 79:
		if prefix == "Y" {
			return true
		}
		if prefix == "Z" {
			return true
		}
	case 80, 81:
		if prefix == "Z" {
			return true
		}
	case 82:
		if prefix == "Z" {
			return true
		}
		if prefix == "a" {
			return true
		}
	case 83:
		if prefix == "a" {
			return true
		}
	case 84:
		if prefix == "a" {
			return true
		}
		if prefix == "b" {
			return true
		}
	case 85:
		if prefix == "b" {
			return true
		}
	case 86:
		if prefix == "b" {
			return true
		}
		if prefix == "c" {
			return true
		}
	case 87, 88:
		if prefix == "c" {
			return true
		}
	case 89:
		if prefix == "c" {
			return true
		}
		if prefix == "d" {
			return true
		}
	case 90:
		if prefix == "d" {
			return true
		}
	case 91:
		if prefix == "d" {
			return true
		}
		if prefix == "e" {
			return true
		}
	case 92, 93:
		if prefix == "e" {
			return true
		}
	case 94:
		if prefix == "e" {
			return true
		}
		if prefix == "f" {
			return true
		}
	case 95:
		if prefix == "f" {
			return true
		}
	case 96:
		if prefix == "f" {
			return true
		}
		if prefix == "g" {
			return true
		}
	case 97, 98:
		if prefix == "g" {
			return true
		}
	case 99:
		if prefix == "g" {
			return true
		}
		if prefix == "h" {
			return true
		}
	case 100:
		if prefix == "h" {
			return true
		}
	case 101:
		if prefix == "h" {
			return true
		}
		if prefix == "i" {
			return true
		}
	case 102, 103:
		if prefix == "i" {
			return true
		}
	case 104:
		if prefix == "i" {
			return true
		}
		if prefix == "j" {
			return true
		}
	case 105:
		if prefix == "j" {
			return true
		}
	case 106:
		if prefix == "j" {
			return true
		}
		if prefix == "k" {
			return true
		}
	case 107, 108:
		if prefix == "k" {
			return true
		}
	case 109:
		if prefix == "k" {
			return true
		}
		if prefix == "m" {
			return true
		}
	case 110:
		if prefix == "m" {
			return true
		}
	case 111:
		if prefix == "m" {
			return true
		}
		if prefix == "n" {
			return true
		}
	case 112, 113:
		if prefix == "n" {
			return true
		}
	case 114:
		if prefix == "n" {
			return true
		}
		if prefix == "o" {
			return true
		}
	case 115:
		if prefix == "o" {
			return true
		}
	case 116:
		if prefix == "o" {
			return true
		}
		if prefix == "p" {
			return true
		}
	case 117, 118:
		if prefix == "p" {
			return true
		}
	case 119:
		if prefix == "p" {
			return true
		}
		if prefix == "q" {
			return true
		}
	case 120:
		if prefix == "q" {
			return true
		}
	case 121:
		if prefix == "q" {
			return true
		}
		if prefix == "r" {
			return true
		}
	case 122, 123:
		if prefix == "r" {
			return true
		}
	case 124:
		if prefix == "r" {
			return true
		}
		if prefix == "s" {
			return true
		}
	case 125:
		if prefix == "s" {
			return true
		}
	case 126:
		if prefix == "s" {
			return true
		}
		if prefix == "t" {
			return true
		}
	case 127, 128:
		if prefix == "t" {
			return true
		}
	case 129:
		if prefix == "t" {
			return true
		}
		if prefix == "u" {
			return true
		}
	case 130:
		if prefix == "u" {
			return true
		}
	case 131:
		if prefix == "u" {
			return true
		}
		if prefix == "v" {
			return true
		}
	case 132, 133:
		if prefix == "v" {
			return true
		}
	case 134:
		if prefix == "v" {
			return true
		}
		if prefix == "w" {
			return true
		}
	case 135:
		if prefix == "w" {
			return true
		}
	case 136:
		if prefix == "w" {
			return true
		}
		if prefix == "x" {
			return true
		}
	case 137, 138:
		if prefix == "x" {
			return true
		}
	case 139:
		if prefix == "x" {
			return true
		}
		if prefix == "y" {
			return true
		}
	case 140:
		if prefix == "y" {
			return true
		}
	case 141:
		if prefix == "y" {
			return true
		}
		if prefix == "z" {
			return true
		}
	case 142, 143:
		if prefix == "z" {
			return true
		}
	case 144:
		if prefix == "z" {
			return true
		}
		if prefix == "2" {
			return true
		}
	}

	if b > 144 {
		if prefix == "2" {
			return true
		}
	}
	return false
}

func TestCanPatternBeSolvedInNet() bool {
	log.Printf("%v", CanPatternBeSolvedInNet("BMog", 26))
	log.Printf("%v", CanPatternBeSolvedInNet("LGWL", 48))
	log.Printf("%v", CanPatternBeSolvedInNet("NPhaux", 52))
	log.Printf("%v", CanPatternBeSolvedInNet("1TPiDev", 0))
	log.Printf("%v", CanPatternBeSolvedInNet("1ThePiachu", 0))
	log.Printf("%v", CanPatternBeSolvedInNet("1Piachu", 0))
	log.Printf("%v", CanPatternBeSolvedInNet("Mvfr8yq8M6yAZwuu6zoVYKpCVQoAv571Qa", 52))
	log.Printf("%v", CanPatternBeSolvedInNet("NATX6zEUNfxfvgVwz8qVnnw3hLhhYXhgQn", 52))
	log.Printf("%v", CanPatternBeSolvedInNet("LhK2kQwiaAvhjWY799cZvMyYwnQAcxkarr", 48))
	log.Printf("%v", CanPatternBeSolvedInNet("xoKDFH4uWpyzxUcCC5jCLFujRKayv3HHcV", 138))
	log.Printf("%v", CanPatternBeSolvedInNet("3EktnHQD7RiAE6uzMj2ZifT9YgRrkSgzQX", 5))
	log.Printf("%v", CanPatternBeSolvedInNet("fF6o8LeDAfswEpMbCW8BqaqmzMWS7TGgew", 95))

	return CanPatternBeSolvedInNet("BMog", 26) == false &&
		CanPatternBeSolvedInNet("LGWL", 48) == false &&
		CanPatternBeSolvedInNet("NPhaux", 52) == false &&
		CanPatternBeSolvedInNet("1TPiDev", 0) &&
		CanPatternBeSolvedInNet("1ThePiachu", 0) &&
		CanPatternBeSolvedInNet("1Piachu", 0) &&
		CanPatternBeSolvedInNet("Mvfr8yq8M6yAZwuu6zoVYKpCVQoAv571Qa", 52) &&
		CanPatternBeSolvedInNet("NATX6zEUNfxfvgVwz8qVnnw3hLhhYXhgQn", 52) &&
		CanPatternBeSolvedInNet("LhK2kQwiaAvhjWY799cZvMyYwnQAcxkarr", 48) &&
		CanPatternBeSolvedInNet("xoKDFH4uWpyzxUcCC5jCLFujRKayv3HHcV", 138) &&
		CanPatternBeSolvedInNet("3EktnHQD7RiAE6uzMj2ZifT9YgRrkSgzQX", 5) &&
		CanPatternBeSolvedInNet("fF6o8LeDAfswEpMbCW8BqaqmzMWS7TGgew", 95)
}

func TestWillNetworkByteGivePrefix() bool {
	return WillNetworkByteGivePrefix(0x00, "1") &&
		WillNetworkByteGivePrefix(5, "3") &&
		WillNetworkByteGivePrefix(48, "L") &&
		WillNetworkByteGivePrefix(52, "M") &&
		WillNetworkByteGivePrefix(52, "N") &&
		WillNetworkByteGivePrefix(95, "f") &&
		WillNetworkByteGivePrefix(97, "g") &&
		WillNetworkByteGivePrefix(105, "j") &&
		WillNetworkByteGivePrefix(111, "m") &&
		WillNetworkByteGivePrefix(111, "n") &&
		WillNetworkByteGivePrefix(125, "s") &&
		WillNetworkByteGivePrefix(127, "t") &&
		WillNetworkByteGivePrefix(128, "t") &&
		//WillNetworkByteGivePrefix(128, "5") && TODO: check if the prefixes are correct on the wiki
		WillNetworkByteGivePrefix(138, "x") &&
		WillNetworkByteGivePrefix(196, "2") &&
		WillNetworkByteGivePrefix(239, "2")
}
