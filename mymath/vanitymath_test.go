package mymath_test

import (
	"testing"

	. "github.com/ThePiachu/Go/mymath"
	"github.com/ThePiachu/Go/mymath/bitelliptic"
)

func TestPrivateKeyToCompressedAddress(t *testing.T) {
	priv := "18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725"
	compressed := "1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs"
	if PrivateKeyToCompressedAddress("0", priv) != compressed {
		t.Error("Invalid compressed address")
	}
	if PrivateKeyToCompressedAddressBytes(0x00, Str2Hex(priv)) != compressed {
		t.Error("Invalid compressed address")
	}

}

func TestPrivateKeyToUncompressedAddress(t *testing.T) {
	priv := "18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725"
	uncompressed := "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM"
	if PrivateKeyToUncompressedAddress("0", priv) != uncompressed {
		t.Error("Invalid uncompressed address")
	}
	if PrivateKeyToUncompressedAddressBytes(0x00, Str2Hex(priv)) != uncompressed {
		t.Error("Invalid uncompressed address")
	}
}

func TestPublicKey(t *testing.T) {
	keys := []string{"0413828ADDF2F27BA6B75856156A295E9E3AD61F2F8788917DFFA17E4DA73D73DACE16B12B2950C39843FB7C9A3A480F6D55DBB9EFD3E70578748EFB9A7B006894",
		"0450863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B23522CD470243453A299FA9E77237716103ABC11A1DF38855ED6F2EE187E9C582BA6"}
	for _, key := range keys {
		a, b := PublicKeyToPointCoordinates(key)
		if a == nil || b == nil {
			t.Errorf("Can't parse the public key %v", key)
			continue
		}

		curve := bitelliptic.S256()
		if !curve.IsOnCurve(a, b) {
			t.Errorf("Public key is not on curve - %v", key)
			continue
		}
	}
}
