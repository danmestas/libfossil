package auth

import "testing"

func TestInviteTokenRoundTrip(t *testing.T) {
	tok := InviteToken{
		URL:      "nats://100.78.32.45:4222/my-repo",
		Login:    "bob",
		Password: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
		Caps:     "oi",
	}
	encoded := tok.Encode()
	if encoded == "" {
		t.Fatal("Encode returned empty string")
	}
	decoded, err := DecodeInviteToken(encoded)
	if err != nil {
		t.Fatalf("DecodeInviteToken: %v", err)
	}
	if decoded.URL != tok.URL {
		t.Errorf("URL = %q, want %q", decoded.URL, tok.URL)
	}
	if decoded.Login != tok.Login {
		t.Errorf("Login = %q, want %q", decoded.Login, tok.Login)
	}
	if decoded.Password != tok.Password {
		t.Errorf("Password = %q, want %q", decoded.Password, tok.Password)
	}
	if decoded.Caps != tok.Caps {
		t.Errorf("Caps = %q, want %q", decoded.Caps, tok.Caps)
	}
}

func TestDecodeInviteTokenInvalid(t *testing.T) {
	_, err := DecodeInviteToken("not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	_, err = DecodeInviteToken("bm90LWpzb24=") // "not-json" in base64
	if err == nil {
		t.Fatal("expected error for non-JSON token")
	}
}
