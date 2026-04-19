package xfer

import "testing"

// TestCardTypeConstants verifies the enum values are correct.
func TestCardTypeConstants(t *testing.T) {
	tests := []struct {
		ct   CardType
		want int
	}{
		{CardFile, 0},
		{CardCFile, 1},
		{CardIGot, 2},
		{CardGimme, 3},
		{CardLogin, 4},
		{CardPush, 5},
		{CardPull, 6},
		{CardCookie, 7},
		{CardClone, 8},
		{CardCloneSeqNo, 9},
		{CardConfig, 10},
		{CardReqConfig, 11},
		{CardPrivate, 12},
		{CardUVFile, 13},
		{CardUVGimme, 14},
		{CardUVIGot, 15},
		{CardPragma, 16},
		{CardError, 17},
		{CardMessage, 18},
		{CardUnknown, 19},
	}
	for _, tt := range tests {
		if int(tt.ct) != tt.want {
			t.Errorf("CardType %d: got %d, want %d", tt.want, int(tt.ct), tt.want)
		}
	}
}

// TestCardInterface verifies that every card struct implements the Card interface
// and returns the correct CardType.
func TestCardInterface(t *testing.T) {
	cards := []struct {
		name string
		card Card
		want CardType
	}{
		{"FileCard", &FileCard{}, CardFile},
		{"CFileCard", &CFileCard{}, CardCFile},
		{"IGotCard", &IGotCard{}, CardIGot},
		{"GimmeCard", &GimmeCard{}, CardGimme},
		{"LoginCard", &LoginCard{}, CardLogin},
		{"PushCard", &PushCard{}, CardPush},
		{"PullCard", &PullCard{}, CardPull},
		{"CookieCard", &CookieCard{}, CardCookie},
		{"CloneCard", &CloneCard{}, CardClone},
		{"CloneSeqNoCard", &CloneSeqNoCard{}, CardCloneSeqNo},
		{"ConfigCard", &ConfigCard{}, CardConfig},
		{"ReqConfigCard", &ReqConfigCard{}, CardReqConfig},
		{"PrivateCard", &PrivateCard{}, CardPrivate},
		{"UVFileCard", &UVFileCard{}, CardUVFile},
		{"UVGimmeCard", &UVGimmeCard{}, CardUVGimme},
		{"UVIGotCard", &UVIGotCard{}, CardUVIGot},
		{"PragmaCard", &PragmaCard{}, CardPragma},
		{"ErrorCard", &ErrorCard{}, CardError},
		{"MessageCard", &MessageCard{}, CardMessage},
		{"UnknownCard", &UnknownCard{}, CardUnknown},
	}
	for _, tt := range cards {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.Type(); got != tt.want {
				t.Errorf("%s.Type() = %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}

// TestCardFieldAccess verifies that card structs hold their field values.
func TestCardFieldAccess(t *testing.T) {
	fc := &FileCard{UUID: "abc123", DeltaSrc: "def456", Content: []byte("data")}
	if fc.UUID != "abc123" || fc.DeltaSrc != "def456" || string(fc.Content) != "data" {
		t.Error("FileCard fields not set correctly")
	}

	igot := &IGotCard{UUID: "abc123", IsPrivate: true}
	if igot.UUID != "abc123" || !igot.IsPrivate {
		t.Error("IGotCard fields not set correctly")
	}

	login := &LoginCard{User: "test user", Nonce: "nonce1", Signature: "sig1"}
	if login.User != "test user" || login.Nonce != "nonce1" || login.Signature != "sig1" {
		t.Error("LoginCard fields not set correctly")
	}

	pragma := &PragmaCard{Name: "sym-pressure", Values: []string{"100"}}
	if pragma.Name != "sym-pressure" || len(pragma.Values) != 1 || pragma.Values[0] != "100" {
		t.Error("PragmaCard fields not set correctly")
	}

	uvig := &UVIGotCard{Name: "file.txt", MTime: 1700000000, Hash: "abc", Size: 42}
	if uvig.MTime != 1700000000 || uvig.Size != 42 {
		t.Error("UVIGotCard fields not set correctly")
	}

	unk := &UnknownCard{Command: "foobar", Args: []string{"a", "b"}}
	if unk.Command != "foobar" || len(unk.Args) != 2 {
		t.Error("UnknownCard fields not set correctly")
	}
}
