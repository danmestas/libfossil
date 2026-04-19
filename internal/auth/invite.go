package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type InviteToken struct {
	URL      string `json:"url"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Caps     string `json:"caps"`
}

func (t InviteToken) Encode() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(fmt.Sprintf("auth.InviteToken.Encode: %v", err))
	}
	return base64.URLEncoding.EncodeToString(b)
}

func DecodeInviteToken(s string) (InviteToken, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return InviteToken{}, fmt.Errorf("auth.DecodeInviteToken: %w", err)
	}
	var t InviteToken
	if err := json.Unmarshal(b, &t); err != nil {
		return InviteToken{}, fmt.Errorf("auth.DecodeInviteToken: %w", err)
	}
	if t.Login == "" {
		return InviteToken{}, fmt.Errorf("auth.DecodeInviteToken: missing login")
	}
	if t.Password == "" {
		return InviteToken{}, fmt.Errorf("auth.DecodeInviteToken: missing password")
	}
	return t, nil
}
