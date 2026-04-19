package delta

import "testing"

func TestApply_NilSource(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil source")
		}
	}()
	Apply(nil, []byte("0\n0;"))
}

func TestCreate_NilSource(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil source")
		}
	}()
	Create(nil, []byte("target"))
}

func TestCreate_NilTarget(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil target")
		}
	}()
	Create([]byte("source"), nil)
}

func TestChecksum_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil data")
		}
	}()
	Checksum(nil)
}
