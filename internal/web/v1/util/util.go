package util

import (
	"fmt"
	"strconv"
	"strings"
)

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
