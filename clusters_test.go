package main

import (
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		size         uint
		clusterShift uint
		tsMask       uint64
		clusterMask  uint64
	}{
		{4, 60, 0x0fffffffffffffff, 0xf000000000000000},
	}

	for _, test := range tests {
		SetClusterBits(test.size)

		if clusterBits != uint(test.size) {
			t.Errorf("%v: expected cluster bits = %v, was %v",
				test.size, test.size, clusterBits)
		}

		if clusterMask != test.clusterMask {
			t.Errorf("%v: expected cluster mask = %x, was %x",
				test.size, test.clusterMask, clusterMask)
		}

		if timestampMask != test.tsMask {
			t.Errorf("%v: expected timestamp mask = %x, was %x",
				test.size, test.tsMask, timestampMask)
		}
	}
}

func TestClusterBits(t *testing.T) {
	SetClusterBits(4)

	tests := []struct {
		input     uint64
		cluster   uint
		timestamp uint64
	}{
		{0, 0, 0},
		{42, 0, 42},
		{0x8088a88f883230f3, 8, 0x088a88f883230f3},
	}

	for _, test := range tests {
		rc := RawCAS(test.input)
		if rc.ClusterID() != test.cluster {
			t.Errorf("Error on cluster ID of %v (%x), exp %v, got %v",
				rc, test.input, test.cluster, rc.ClusterID())
		}
		if rc.Timestamp() != test.timestamp {
			t.Errorf("Error on timestamp ID of %v (%x), exp %v, got %v",
				rc, test.input, test.timestamp, rc.Timestamp())
		}
	}
}
