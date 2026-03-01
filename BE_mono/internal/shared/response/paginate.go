package response

type PageDto[T any] struct {
	Items      []T      `json:"items"`
	Pagination PageMeta `json:"pagination"`
}

type PageMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
