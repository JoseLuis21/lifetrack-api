package shared

// FilterMap generic and non-type safe filter map
type FilterMap map[string]string

// CategoryCriteria represents a filter struct to fetch fine-grained categories
type CategoryCriteria struct {
	User    string `json:"user"`
	Query   string `json:"query"`
	OrderBy string `json:"order_by"`
}

// ActivityCriteria represents a filter struct to fetch fine-grained activities
type ActivityCriteria struct {
	Category string `json:"category"`
	Query    string `json:"query"`
	OrderBy  string `json:"order_by"`
}

// ActivityCriteria represents a filter struct to fetch fine-grained occurrences
type OccurrenceCriteria struct {
	Activity string `json:"activity"`
	Query    string `json:"query"`
	OrderBy  string `json:"order_by"`
}
