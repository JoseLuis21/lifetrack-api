package categoryhandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/neutrinocorp/life-track-api/internal/application/category"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
)

type Get struct {
	q      *category.GetQuery
	router *mux.Router
}

// NewGet creates a Get handler with routing
func NewGet(q *category.GetQuery, r *mux.Router) *Get {
	h := &Get{
		q:      q,
		router: r,
	}
	h.mapRoutes()

	return h
}

func (c *Get) mapRoutes() {
	c.router.Path("/category/{id}").Methods(http.MethodGet).HandlerFunc(c.Handler)
}

func (c Get) GetRouter() *mux.Router {
	return c.router
}

func (c *Get) Handler(w http.ResponseWriter, r *http.Request) {
	cat, err := c.q.Query(context.Background(), mux.Vars(r)["id"])
	if err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(cat)
}
