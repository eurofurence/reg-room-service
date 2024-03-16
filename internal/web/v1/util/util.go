package util

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// NewStrictJSONDecoder constructs a new JSON decoder with strict settings
func NewStrictJSONDecoder(r io.Reader) *json.Decoder {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	return dec
}

// ParseMemberIDs is a helper function to parse member IDs for groups and rooms.
func ParseMemberIDs(ids string) ([]string, error) {
	var res []string
	if ids != "" {
		res = make([]string, 0)

		// we expect IDs to be provided as a comma separated list
		// ids must be numeric. If any ID is invalid we want to return an error
		splitIDs := strings.Split(ids, ",")
		for _, memberID := range splitIDs {
			if !IsNumeric(memberID) {
				return nil, fmt.Errorf("member ids must be numeric and valid. Invalid member id: %s", memberID)
			}
			res = append(res, memberID)
		}
	}

	return res, nil
}

// IsNumeric is a helper function to determine if a
// string is a number.
func IsNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	signed | unsigned
}

func ParseInt[E signed](s string) (E, error) {
	res, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	casted := E(res)
	if int(casted) != res {
		return E(0), fmt.Errorf("could not cast %q to target type %T", s, casted)
	}

	return E(res), nil
}

// ParseOptionalBool will return false if the input string is empty,
//
// If the sting is not empty and the parsed value is invalid
// it will return an error instead.
//
// Otherwise, it will return the correctly parsed bool value.
func ParseOptionalBool(s string) (bool, error) {
	if s == "" {
		return false, nil
	}

	return strconv.ParseBool(s)
}

func ParseUInt[E unsigned](s string) (E, error) {
	res, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 0, err
	}

	casted := E(res)
	if uint(casted) != uint(res) {
		return E(0), fmt.Errorf("could not cast %q to target type %T", s, casted)
	}

	return E(res), nil
}

type Float interface {
	~float32 | ~float64
}

func ParseFloat[E Float](s string) (E, error) {
	var res E

	bitSize := 32

	switch any(res).(type) {
	case float64:
		bitSize = 64
	}

	parsed, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return E(0), err
	}

	return E(parsed), nil
}
