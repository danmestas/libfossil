package diff

import "testing"

func FuzzUnified(f *testing.F) {
	f.Add([]byte("a\nb\n"), []byte("a\nc\n"))
	f.Add([]byte(""), []byte("hello\n"))
	f.Add([]byte("hello\n"), []byte(""))
	f.Add([]byte("same\n"), []byte("same\n"))
	f.Add([]byte("\x00binary"), []byte("text\n"))
	f.Add([]byte("a\r\nb\r\n"), []byte("a\nb\n"))

	f.Fuzz(func(t *testing.T, a, b []byte) {
		// Must not panic (negative ContextLines is caught, but fuzz uses 3).
		_ = Unified(a, b, Options{ContextLines: 3})

		stat := Stat(a, b)
		if stat.Binary {
			return
		}
		if stat.Insertions < 0 {
			t.Fatalf("negative insertions: %d", stat.Insertions)
		}
		if stat.Deletions < 0 {
			t.Fatalf("negative deletions: %d", stat.Deletions)
		}
	})
}
