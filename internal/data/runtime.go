package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Error for UnmarshalJSON() when it is unable to parse the JSON string.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// Custom type on Runtime value.
type Runtime int32

// MarshalJSON specified JSON encoding approach.
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

// UnmarshalJSON specified JSON decoding approach.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// remove double-quotes from string "<runtime> mins".
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// check format.
	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// parse int32.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// assign to the receiver.
	*r = Runtime(i)
	return nil
}
