package idkit

import (
	"crypto/rand"
	"fmt"

	"github.com/denisbrodbeck/machineid"
)

const (
	entropyZero byte = iota
	entropyQuick
	entropyMachPid
)

func readMachineID() []byte {
	id, err := machineid.ID()
	if err != nil {
		panic(err)
	}

	return []byte(id)
}

// randInt generates a random uint32
func randInt() uint32 {
	b := make([]byte, 3)
	if _, err := rand.Reader.Read(b); err != nil {
		panic(fmt.Errorf("xid: cannot generate random number: %v;", err))
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}
