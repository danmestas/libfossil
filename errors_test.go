package libfossil

import (
	"errors"
	"testing"
)

func TestFslErrorCode(t *testing.T) {
	err := &FslError{Code: RCNotARepo, Msg: "not a repo"}
	if err.Code != RCNotARepo {
		t.Fatalf("Code = %d, want %d", err.Code, RCNotARepo)
	}
}

func TestFslErrorImplementsError(t *testing.T) {
	var err error = &FslError{Code: RCDB, Msg: "db error"}
	got := err.Error()
	want := "fossil(db): db error"
	if got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestFslErrorUnwrap(t *testing.T) {
	inner := errors.New("sqlite: constraint")
	err := &FslError{Code: RCDB, Msg: "insert failed", Cause: inner}
	if !errors.Is(err, inner) {
		t.Fatal("FslError should unwrap to its Cause")
	}
}

func TestRCCodeString(t *testing.T) {
	tests := []struct {
		code RC
		want string
	}{
		{RCOK, "ok"},
		{RCError, "error"},
		{RCOOM, "oom"},
		{RCMisuse, "misuse"},
		{RCRange, "range"},
		{RCAccess, "access"},
		{RCIO, "io"},
		{RCNotFound, "not_found"},
		{RCAlreadyExists, "already_exists"},
		{RCConsistency, "consistency"},
		{RCNotARepo, "not_a_repo"},
		{RCDB, "db"},
		{RCChecksumMismatch, "checksum_mismatch"},
	}
	for _, tt := range tests {
		if got := tt.code.String(); got != tt.want {
			t.Errorf("RC(%d).String() = %q, want %q", tt.code, got, tt.want)
		}
	}
}
