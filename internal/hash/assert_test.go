package hash

import "testing"

func TestSHA1_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil data")
		}
	}()
	SHA1(nil)
}

func TestSHA3_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil data")
		}
	}()
	SHA3(nil)
}

func TestHashSize_Empty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty hashType")
		}
	}()
	HashSize("")
}

func TestIsValidHash_Empty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty hash")
		}
	}()
	IsValidHash("")
}
