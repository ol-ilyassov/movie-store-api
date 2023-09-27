package data

import (
	"strings"

	"movie.api/internal/validator"
)

// Filters
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string // holds supported sort values.
}

func ValidateFilters(v *validator.Validator, f Filters) {
	// validate for sensible values.
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

// sortColumn returns column if appropriate and trims prefix.
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	// could be implemented as return error.
	// by implementation, it is already handled by ValidateFilters().
	// this is a sensible failsafe to help stop a SQL injection attack occurring.
	panic("unsafe sort parameter: " + f.Sort)
}

// sortDirection returns sort order keyword depending on sort value prefix.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
