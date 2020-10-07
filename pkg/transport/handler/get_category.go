package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/query"
)

type GetCategory struct {
	q      *query.GetCategory
	router *mux.Router
}

// NewGetCategory creates a get category handler with routing
func NewGetCategory(q *query.GetCategory, r *mux.Router) *GetCategory {
	h := &GetCategory{
		q:      q,
		router: r,
	}
	h.mapRoutes()

	return h
}

func (c *GetCategory) mapRoutes() {
	c.router.Path("/live/category/{id}").Methods(http.MethodGet).HandlerFunc(c.Handler)
}

func (c GetCategory) GetRouter() *mux.Router {
	return c.router
}

func (c *GetCategory) Handler(w http.ResponseWriter, r *http.Request) {
	category, err := c.q.Query(context.Background(), mux.Vars(r)["id"])
	if err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(category)
}
