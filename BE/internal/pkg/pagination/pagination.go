package pagination

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

// Pagination carries page and limit parameters with sensible defaults.
type Pagination struct {
	Page  int
	Limit int
}

// New creates a Pagination with defaults applied and max-limit enforcement.
func New(page, limit int) Pagination {
	if page <= 0 {
		page = defaultPage
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return Pagination{Page: page, Limit: limit}
}

// Offset returns the SQL offset for the current page.
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}
