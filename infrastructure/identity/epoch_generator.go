package identity

import "github.com/smartystreets/clock"

type EpochGenerator struct {
	clock *clock.Clock
	last  uint64
}

func NewEpochGenerator() *EpochGenerator {
	return &EpochGenerator{}
}

func (this *EpochGenerator) Next() uint64 {
	if id := this.next(); this.last >= id {
		this.last++
		return this.last
	} else {
		this.last = id
		return id
	}
}
func (this *EpochGenerator) next() uint64 {
	return uint64(this.clock.UTCNow().Unix())
}
