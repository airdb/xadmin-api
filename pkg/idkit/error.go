package idkit

const (
	// ErrInvalidID is returned when trying to unmarshal an invalid ID.
	ErrInvalidID strErr = "idkit: invalid ID"
)

// strErr allows declaring errors as constants.
type strErr string

func (err strErr) Error() string { return string(err) }
