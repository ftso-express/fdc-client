package attestation

type Manager struct {
	Rounds         map[uint64]Round
	Timestamps     chan uint64
	Requests       chan uint64
	BitVotes       chan uint64
	RoundInCollect uint64
}
