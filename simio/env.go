package simio

// Env bundles all I/O abstractions needed for deterministic simulation.
// Pass *Env through call chains instead of threading individual interfaces.
type Env struct {
	Clock   Clock
	Rand    Rand
	Storage Storage
}

// RealEnv returns an Env wired to real system I/O (production use).
func RealEnv() *Env {
	return &Env{
		Clock:   RealClock{},
		Rand:    CryptoRand{},
		Storage: OSStorage{},
	}
}

// SimEnv returns an Env wired to deterministic implementations
// controlled by the given seed.
func SimEnv(seed int64) *Env {
	return &Env{
		Clock:   NewSimClock(),
		Rand:    NewSeededRand(seed),
		Storage: NewMemStorage(),
	}
}
