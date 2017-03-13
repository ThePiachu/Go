package mymath

// Copyright 2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func AESEncrypt(key string, plaintext string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return ciphertext, nil
}

func AESDecrypt(key string, ciphertext []byte) (string, error) {
	resp, err := AESDecryptBytes([]byte(key), ciphertext)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func AESDecryptBytes(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}

func AESEncryptCBC(key string, plaintext string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	paddedText := PadCBC([]byte(plaintext))
	ciphertext := make([]byte, aes.BlockSize+len(paddedText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(ciphertext[aes.BlockSize:], paddedText)

	return ciphertext, nil
}

func AESDecryptBytesCBC(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)
	return UnpadCBC(ciphertext)
}

func AESDecryptCBC(key string, ciphertext []byte) (string, error) {
	resp, err := AESDecryptBytesCBC([]byte(key), ciphertext)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func TrimISO10126Padding(b []byte) []byte {
	if len(b) < 1 {
		return nil
	}
	padLength := int(b[len(b)-1])
	if len(b) < padLength {
		return nil
	}
	return b[:len(b)-padLength]
}

func PadCBC(b []byte) []byte {
	l := len(b)
	padLen := 16 - (l % 16)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(b, padText...)
}

func UnpadCBC(b []byte) ([]byte, error) {
	l := len(b)
	padLen := int(b[l-1])
	if padLen >= l || padLen > 16 {
		return nil, errors.New("Invalid padding size")
	}
	return b[:l-padLen], nil
}
