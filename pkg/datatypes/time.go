package datatypes

import (
	"bytes"
	"time"
)

type Time struct {
	*time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format("\"" + time.RFC3339 + "\"")), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(data []byte) (err error) {

	// by convention, unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
	if bytes.Equal(data, []byte("null")) ||
		bytes.Equal(data, []byte(`""`)) {
		return nil
	}

	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse("\""+time.RFC3339+"\"", string(data))
	*t = Time{&tt}
	return
}
