package db

import "fmt"

// Paginate does the API-wide processing of pagination controls:
// maximum limit is 100, default limit is 30, default offset is zero.
func Paginate(inLimit, inOffset *uint) string {
	limit := uint(30)
	if inLimit != nil {
		limit = *inLimit
	}
	if limit > 100 {
		limit = uint(100)
	}
	offset := uint(0)
	if inOffset != nil {
		offset = *inOffset
	}
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}
