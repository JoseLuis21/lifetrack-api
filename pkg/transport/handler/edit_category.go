package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
)

type EditCategory struct {
	cmd    *command.EditCategoryHandler
	router *mux.Router
}

// NewEditCategory creates an edit category handler with routing
func NewEditCategory(cmd *command.EditCategoryHandler, r *mux.Router) *EditCategory {
	h := &EditCategory{
		cmd:    cmd,
		router: r,
	}
	h.mapRoute()

	return h
}

func (c *EditCategory) mapRoute() {
	c.router.Path("/category/{id}").Methods(http.MethodPut, http.MethodPatch).HandlerFunc(c.Handler)
}

func (c EditCategory) GetRouter() *mux.Router {
	return c.router
}

func (c EditCategory) Handler(w http.ResponseWriter, r *http.Request) {
	if err := c.cmd.Invoke(command.EditCategory{
		Ctx:         r.Context(),
		ID:          mux.Vars(r)["id"],
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
		Theme:       r.PostFormValue("theme"),
	}); err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(httputil.Response{
		Message: "successfully updated category",
		Code:    http.StatusOK,
	})
}
