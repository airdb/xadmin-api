package idkit

import (
	cryptorand "crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"hash/crc32"
	"io"
	mathrand "math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/oklog/ulid/v2"
)

type Id struct {
	ulid.ULID
}

const (
	rawLen             = 16
	version       byte = 0
	entropyMethod      = entropyMachPid
	entropyOffset      = 6
)

var (
	// objectIDCounter is atomically incremented when generating a new ObjectId
	// using NewObjectId() function. It's used as a counter part of an id.
	// This id is initialized with a random value.
	objectIDCounter = randInt()

	// machineId stores machine id generated once and used in subsequent calls
	// to NewObjectId function.
	machineId = readMachineID()

	// pid stores the current process id
	pid = os.Getpid()

	crc32q = crc32.New(crc32.MakeTable(crc32.Koopman))

	nilID Id
)

// New generates a globally unique ID
func New() Id {
	now := time.Now().UTC()
	return NewWithTime(now)
}

// NewWithTime generates a globally unique ID with the passed in time
func NewWithTime(t time.Time) Id {
	return Id{
		ulid.MustNew(ulid.Timestamp(t), newMachPidReader()),
	}
}

// FromString reads an ID from its string representation
func FromString(s string) (Id, error) {
	id := &Id{}
	err := id.UnmarshalText([]byte(s))
	return *id, err
}

// MustFromString reads an ID from its string representation
func MustFromString(s string) Id {
	id := &Id{}
	err := id.UnmarshalText([]byte(s))
	if err != nil {
		panic(err)
	}
	return *id
}

// FromBytes convert the byte array representation of `ID` back to `ID`
func FromBytes(b []byte) (Id, error) {
	var id Id
	if len(b) != rawLen {
		return id, ErrInvalidID
	}
	copy(id.ULID[:], b)
	return id, nil
}

// NilID returns a zero value for `ulid.ID`.
func NilID() Id {
	return nilID
}

// IsNil Returns true if this is a "nil" ID
func (id Id) IsNil() bool {
	return id == nilID
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id Id) Time() time.Time {
	// First 4 bytes of ObjectId is 32-bit big-endian seconds from epoch.
	secs := int64(id.ULID.Time())
	return time.UnixMilli(secs)
}

// Machine returns the 3-byte machine id part of the id.
// It's a runtime error to call this method with an invalid id.
func (id Id) Machine() []byte {
	switch id.ULID[entropyOffset+1] {
	case entropyMachPid:
		return id.ULID[entropyOffset+2 : entropyOffset+5]
	default:
		return nil
	}
}

// Pid returns the process id part of the id.
// It's a runtime error to call this method with an invalid id.
func (id Id) Pid() uint16 {
	switch id.ULID[entropyOffset+1] {
	case entropyMachPid:
		return binary.BigEndian.Uint16(id.ULID[entropyOffset+5 : entropyOffset+7])
	default:
		return 0
	}
}

// Counter returns the incrementing value part of the id.
// It's a runtime error to call this method with an invalid id.
func (id Id) Counter() int32 {
	switch id.ULID[entropyOffset+1] {
	case entropyMachPid:
		b := id.ULID[entropyOffset+7 : entropyOffset+10]
		// Counter is stored as big-endian 3-byte value
		return int32(uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]))
	default:
		return 0
	}
}

// Value implements the driver.Valuer interface.
func (id Id) Value() (driver.Value, error) {
	if id.IsNil() {
		return nil, nil
	}
	b, err := id.MarshalText()
	return string(b), err
}

type zeroReader struct{}

func newZeroReader(t time.Time) io.Reader {
	return zeroReader{}
}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

type quickReader struct {
	entropy io.Reader
}

func newQuckReader(t time.Time) io.Reader {
	entropy := cryptorand.Reader
	seed := t.UnixNano()
	source := mathrand.NewSource(seed)
	entropy = mathrand.New(source)

	return &quickReader{entropy}
}

func (r quickReader) Read(p []byte) (int, error) {
	return r.entropy.Read(p)
}

type machPidReader struct {
}

func newMachPidReader() io.Reader {
	return &machPidReader{}
}

func (r machPidReader) Read(p []byte) (int, error) {
	// begin from offset 6, and remian 10 bytes to add custome data.
	p[0] = version
	p[1] = entropyMethod
	// Machine, first 3 bytes of md5(hostname)
	p[2] = machineId[0]
	p[3] = machineId[1]
	p[4] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	p[5] = byte(pid >> 8)
	p[6] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	p[7] = byte(i >> 16)
	p[8] = byte(i >> 8)
	p[9] = byte(i)

	// p[0] = crc32q.Sum(p)[0]&0xF0 | version&0x0F

	return 10, nil
}
