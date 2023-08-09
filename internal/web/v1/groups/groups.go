package groups

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eurofurence/reg-room-service/internal/controller"
)

// Empty defines a type which is used for empty responses
type Empty struct{}

// Handler implements methods, which satisfy the endpoint format
// in the `common` package
type Handler struct {
	ctrl controller.Controller
}

func parseGroupMemberIDs(ids string) ([]string, error) {
	var res []string
	if ids != "" {
		res = make([]string, 0)

		// we expect IDs to be provided as a comma separated list
		// ids must be numeric. If any ID is invalid we want to return an error
		splitIDs := strings.Split(ids, ",")
		for _, memberID := range splitIDs {
			if !isNumeric(memberID) {
				return nil, fmt.Errorf("member ids must be numeric and valid. Invalid member id: %s", memberID)
			}
			res = append(res, memberID)
		}
	}

	return res, nil
}

func isNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
