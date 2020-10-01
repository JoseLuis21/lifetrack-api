package shared

// FilterMap generic and non-type safe filter map
type FilterMap map[string]string

// CategoryCriteria represents a filter struct to fetch fine-grained categories
type CategoryCriteria struct {
	User  string `json:"user"`
	Title string `json:"title"`
	Query string `json:"query"`
}
