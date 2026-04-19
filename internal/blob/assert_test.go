package blob

import "testing"

func TestStore_NilQuerier(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil querier")
		}
	}()
	Store(nil, []byte("content"))
}

func TestStoreDelta_NilQuerier(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil querier")
		}
	}()
	StoreDelta(nil, []byte("content"), 1)
}

func TestStorePhantom_NilQuerier(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil querier")
		}
	}()
	StorePhantom(nil, "abc123")
}

func TestLoad_NilQuerier(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil querier")
		}
	}()
	Load(nil, 1)
}

func TestExists_NilQuerier(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil querier")
		}
	}()
	Exists(nil, "abc123")
}

func TestCompress_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil data")
		}
	}()
	Compress(nil)
}

func TestDecompress_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil data")
		}
	}()
	Decompress(nil)
}
