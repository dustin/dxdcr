package main

import (
	"fmt"
)

type RawCAS uint64

var (
	clusterBits   uint
	clusterMask   uint64
	timestampMask uint64
	clusterShift  uint
)

// Define number of cluster bits
func SetClusterBits(to uint) {
	clusterBits = clusterRangeCheck(4, to)
	clusterMask = (((1 << clusterBits) - 1) << (64 - clusterBits))
	timestampMask = ^clusterMask
	clusterShift = 64 - clusterBits
}

func init() {
	SetClusterBits(4)
}

func clusterRangeCheck(bits, to uint) uint {
	if to >= 1<<4 {
		panic("bits out of range")
	}

	return to
}

// Timestamp is the local sequence identifier.
func (r RawCAS) Timestamp() uint64 {
	return uint64(r) & timestampMask
}

// ClusterID Identifies which cluster originated.
func (r RawCAS) ClusterID() uint {
	return uint(uint64(r) >> clusterShift)
}

func (r RawCAS) String() string {
	return fmt.Sprintf("{cid: %x, ts: %x}", r.ClusterID(), r.Timestamp())
}

func (r *RawCAS) SetTimestamp(ts uint64) {
	*r = (^RawCAS(timestampMask))&*r | RawCAS(ts&clusterMask)
}

func (r *RawCAS) SetClusterID(to uint) {
	*r = RawCAS(timestampMask)&*r | RawCAS(clusterRangeCheck(clusterBits, to))
}
