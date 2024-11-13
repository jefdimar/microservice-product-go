package validation

type QueryParams struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  string
}

func ValidateQueryParams(params *QueryParams) error {
	ve := &ValidationErrors{}

	if params.Page < 1 {
		ve.AddError("page", "Page must be greater than 0")
	}

	if params.PageSize < 1 || params.PageSize > 100 {
		ve.AddError("pageSize", "Page size must be between 1 and 100")
	}

	validSortFields := map[string]bool{"name": true, "price": true, "created_at": true}
	if params.SortDir != "" && params.SortDir != "asc" && params.SortDir != "desc" {
		ve.AddError("sortDir", "Invalid sort direction. Must be either asc or desc.")
	}

	if params.SortBy != "" && !validSortFields[params.SortBy] {
		ve.AddError("sortBy", "Invalid sort field.")
	}

	if ve.HasErrors() {
		return ve
	}

	return nil
}
