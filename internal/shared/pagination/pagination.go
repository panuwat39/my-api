package pagination

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type Query struct {
	Page  int
	Limit int
}

type Metadata struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type Response struct {
	Items any      `json:"items"`
	Meta  Metadata `json:"meta"`
}

func NewQuery(page, limit int) Query {
	if page < 1 {
		page = DefaultPage
	}

	if limit < 1 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Query{
		Page:  page,
		Limit: limit,
	}
}

func (q Query) Skip() int64 {
	return int64((q.Page - 1) * q.Limit)
}

func NewMetadata(
	page int,
	limit int,
	total int64,
) Metadata {
	totalPages := int64(0)

	if total > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}

	return Metadata{
		Page:       int64(page),
		Limit:      int64(limit),
		Total:      total,
		TotalPages: totalPages,
	}
}
