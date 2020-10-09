package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
)

type AddCategory struct {
	cmd    *command.AddCategoryHandler
	router *mux.Router
}

// NewAddCategory creates an add category handler with routing
func NewAddCategory(cmd *command.AddCategoryHandler, r *mux.Router) *AddCategory {
	h := &AddCategory{
		cmd:    cmd,
		router: r,
	}
	h.mapRoute()

	return h
}

func (c *AddCategory) mapRoute() {
	c.router.Path("/category").Methods(http.MethodPost).HandlerFunc(c.Handler)
}

func (c AddCategory) GetRouter() *mux.Router {
	return c.router
}

func (c AddCategory) Handler(w http.ResponseWriter, r *http.Request) {
	if err := c.cmd.Invoke(command.AddCategory{
		Ctx:         r.Context(),
		Title:       r.PostFormValue("title"),
		User:        r.PostFormValue("user"),
		Description: r.PostFormValue("description"),
	}); err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(httputil.Response{
		Message: "successfully created category",
		Code:    http.StatusOK,
	})
}
