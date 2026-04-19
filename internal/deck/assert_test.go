package deck

import "testing"

func TestParseNilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Parse(nil) did not panic")
		}
	}()
	Parse(nil)
}

func TestMarshalNilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("(*Deck)(nil).Marshal() did not panic")
		}
	}()
	var d *Deck
	d.Marshal()
}
