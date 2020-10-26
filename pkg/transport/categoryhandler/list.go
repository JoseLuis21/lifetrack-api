package categoryhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/neutrinocorp/life-track-api/internal/application/category"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
)

type List struct {
	q      *category.ListQuery
	router *mux.Router
}

// NewList creates a new List handler with routing
func NewList(q *category.ListQuery, r *mux.Router) *List {
	h := &List{
		q:      q,
		router: r,
	}
	h.mapRoutes()

	return h
}

func (c *List) mapRoutes() {
	c.router.StrictSlash(true).Path("/category").Methods(http.MethodGet).HandlerFunc(c.Handler)
}

func (c List) GetRouter() *mux.Router {
	return c.router
}

func (c *List) Handler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil {
		limit = 100
	}
	categories, nextPage, err := c.q.Query(r.Context(), category.Filter{
		UserID:  r.URL.Query().Get("user"),
		Name:    r.URL.Query().Get("name"),
		Keyword: r.URL.Query().Get("query"),
		Limit:   limit,
		Token:   r.URL.Query().Get("next_page"),
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
