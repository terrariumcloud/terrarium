package restapi

import (
	"errors"
	"net/url"
	"strconv"
)

// ExtractLimitAndOffset is a convience method to extract pagination and limit values passed
// to a handler that supports pages.
func ExtractLimitAndOffset(qs url.Values) (int, int, error) {
	var limit int = 10
	var offset int = 0
	if qs.Has("limit") {
		// If we have a limit value set in QS attempt to convert to int
		i, err := strconv.Atoi(qs.Get("limit"))
		if err != nil {
			return 0, 0, errors.New("limit is not an integer")
		}
		limit = i
	}
	if qs.Has("offset") {
		i, err := strconv.Atoi(qs.Get("offset"))
		if err != nil {
			return 0, 0, errors.New("offset is not an integer")
		}
		offset = i
	}
	return limit, offset, nil
}
