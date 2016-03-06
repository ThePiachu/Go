package mymath_test

import (
	"testing"
	. "github.com/ThePiachu/Go/mymath"
)

func TestPrivateKeyToCompressedAddress(t *testing.T) {
	priv:="18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725"
	compressed:="1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs"
	if PrivateKeyToCompressedAddress("0", priv)!=compressed {
		t.Error("Invalid compressed address")
	}
	if PrivateKeyToCompressedAddressBytes(0x00, Str2Hex(priv))!=compressed {
		t.Error("Invalid compressed address")
	}

}

func TestPrivateKeyToUncompressedAddress(t *testing.T) {
	priv:="18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725"
	uncompressed:="16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM"
	if PrivateKeyToUncompressedAddress("0", priv)!=uncompressed {
		t.Error("Invalid uncompressed address")
	}
	if PrivateKeyToUncompressedAddressBytes(0x00, Str2Hex(priv))!=uncompressed {
		t.Error("Invalid uncompressed address")
	}
}