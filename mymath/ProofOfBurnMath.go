package mymath

// Copyright 2012 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"errors"
)

func GenerateProofOfBurnAddress(addressRoot string, padding byte) (string, error) {
	if IsBase58StringValid(addressRoot) == false {
		return "", errors.New("Invalid character in input")
	}
	if IsBase58StringValid(string(padding)) == false {
		return "", errors.New("Invalid character in input")
	}

	answer := addressRoot
	hash := Base582Hex(answer)
	if len(hash) > 26 {
		return "", errors.New("Input too long")
	}
	if len(hash) != 25 {
		for len(hash) < 25 {
			answer = answer + string(padding)
			hash = Base582Hex(answer)
		}
	}
	checksum := DoubleSHA(hash[0:21])

	//hash checksum
	hash[21] = checksum[0]
	hash[22] = checksum[1]
	hash[23] = checksum[2]
	hash[24] = checksum[3]

	return Hex2Base58String(hash), nil
}
