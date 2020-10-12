package categoryhandler

import (
	"encoding/json"
	"net/http"

	"github.com/neutrinocorp/life-track-api/internal/application/category"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
)

type Remove struct {
	cmd *category.RemoveHandler
	r   *mux.Router
}

// NewRemove creates a Remove handler with routing
func NewRemove(cmd *category.RemoveHandler, r *mux.Router) *Remove {
	h := &Remove{
		cmd: cmd,
		r:   r,
	}
	h.mapRoutes()

	return h
}

func (c *Remove) mapRoutes() {
	c.r.Path("/category/{id}").Methods(http.MethodDelete).HandlerFunc(c.Handler)
}

func (c Remove) GetRouter() *mux.Router {
	return c.r
}

func (c Remove) Handler(w http.ResponseWriter, r *http.Request) {
	if err := c.cmd.Invoke(category.Remove{
		Ctx: r.Context(),
		ID:  mux.Vars(r)["id"],
	}); err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(httputil.Response{
		Message: "successfully removed category",
		Code:    http.StatusOK,
	})
}
