package libfossil

import "testing"

func TestFslIDIsInt64(t *testing.T) {
	var id FslID = -1
	if id != -1 {
		t.Fatal("FslID should support negative values")
	}
	var big FslID = 1 << 33
	if big <= 0 {
		t.Fatal("FslID should support values > int32 max")
	}
}

func TestFslSizePhantom(t *testing.T) {
	if PhantomSize != -1 {
		t.Fatalf("PhantomSize = %d, want -1", PhantomSize)
	}
	var s FslSize = PhantomSize
	if s >= 0 {
		t.Fatal("FslSize should be able to represent -1 for phantoms")
	}
}

func TestFossilAppID(t *testing.T) {
	if FossilApplicationID != 252006673 {
		t.Fatalf("FossilApplicationID = %d, want 252006673", FossilApplicationID)
	}
}

func TestCreateOptsDefaults(t *testing.T) {
	opts := CreateOpts{}
	if opts.User != "" {
		t.Error("default User should be empty")
	}
}

func TestSyncOptsDefaults(t *testing.T) {
	opts := SyncOpts{}
	if opts.Push || opts.Pull {
		t.Error("default Push/Pull should be false")
	}
}

func TestLogOptsLimit(t *testing.T) {
	opts := LogOpts{Limit: 10}
	if opts.Limit != 10 {
		t.Errorf("got %d, want 10", opts.Limit)
	}
}
