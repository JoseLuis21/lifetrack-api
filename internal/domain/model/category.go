package model

// Category read model
type Category struct {
	ID          string `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	User        string `json:"user"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	Active      bool   `json:"active"`
}
