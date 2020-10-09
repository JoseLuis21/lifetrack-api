package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
)

type RemoveCategory struct {
	cmd *command.RemoveCategoryHandler
	r   *mux.Router
}

// NewRemoveCategory creates an remove category state handler with routing
func NewRemoveCategory(cmd *command.RemoveCategoryHandler, r *mux.Router) *RemoveCategory {
	h := &RemoveCategory{
		cmd: cmd,
		r:   r,
	}
	h.mapRoutes()

	return h
}

func (c *RemoveCategory) mapRoutes() {
	c.r.Path("/category/{id}").Methods(http.MethodDelete).HandlerFunc(c.Handler)
}

func (c RemoveCategory) GetRouter() *mux.Router {
	return c.r
}

func (c RemoveCategory) Handler(w http.ResponseWriter, r *http.Request) {
	if err := c.cmd.Invoke(command.RemoveCategory{
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
