package auth

import (
	"strings"
	"time"
)

type User struct {
	UID     int
	Login   string
	Cap     string
	CExpire time.Time // zero value = no expiry
	Info    string
	MTime   time.Time
}

func HasCapability(caps string, required byte) bool {
	return strings.IndexByte(caps, required) >= 0
}

func CanPush(caps string) bool        { return HasCapability(caps, 'i') }
func CanPull(caps string) bool        { return HasCapability(caps, 'o') }
func CanClone(caps string) bool       { return HasCapability(caps, 'g') }
func CanPushUV(caps string) bool      { return HasCapability(caps, 'y') }
func CanSyncPrivate(caps string) bool { return HasCapability(caps, 'x') }
