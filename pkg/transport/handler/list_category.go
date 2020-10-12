package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/query"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
)

type ListCategory struct {
	q      *query.ListCategories
	router *mux.Router
}

// NewListCategory creates a new list category handler with routing
func NewListCategory(q *query.ListCategories, r *mux.Router) *ListCategory {
	h := &ListCategory{
		q:      q,
		router: r,
	}
	h.mapRoutes()

	return h
}

func (c *ListCategory) mapRoutes() {
	c.router.StrictSlash(true).Path("/category").Methods(http.MethodGet).HandlerFunc(c.Handler)
}

func (c ListCategory) GetRouter() *mux.Router {
	return c.router
}

func (c *ListCategory) Handler(w http.ResponseWriter, r *http.Request) {
	categories, nextPage, err := c.q.Query(r.Context(), r.URL.Query().Get("next_page"), r.URL.Query().Get("page_size"), map[string]string{
		"user":  r.URL.Query().Get("user"),
		"query": r.URL.Query().Get("query"),
		"order": r.URL.Query().Get("order"),
	})
	if err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(struct {
		Categories []*model.Category `json:"categories"`
		TotalItems int               `json:"total_items"`
		NextPage   string            `json:"next_page"`
	}{
		Categories: categories,
		TotalItems: len(categories),
		NextPage:   nextPage,
	})
}
