package utils

type PaginationClauses struct {
	Offset int
	Limit  int
}

func OffsetAndLimit(page, size int) PaginationClauses {
	if page < 1 {
		page = 1
	}

	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	return PaginationClauses{
		Offset: offset,
		Limit:  size,
	}
}
