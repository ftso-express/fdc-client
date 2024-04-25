package attestation

type Manager struct {
	Rounds         map[uint64]Round
	BlockTimestamp uint64
}
