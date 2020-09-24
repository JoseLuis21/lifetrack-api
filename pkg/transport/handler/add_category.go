package handler

import (
	"encoding/json"
	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
	"net/http"
)

type AddCategory struct {
	cmd    *command.AddCategoryHandler
	router *mux.Router
}

// NewAddCategory creates a add category handler with routing
func NewAddCategory(cmd *command.AddCategoryHandler) *AddCategory {
	h := &AddCategory{
		cmd:    cmd,
		router: mux.NewRouter(),
	}
	h.mapRoute()

	return h
}

func (c *AddCategory) mapRoute() {
	c.router.Path("/live/category").Methods(http.MethodPost).HandlerFunc(c.Handler)
}

func (c AddCategory) GetRouter() *mux.Router {
	return c.router
}

func (c AddCategory) Handler(w http.ResponseWriter, r *http.Request) {
	err := c.cmd.Invoke(command.AddCategory{
		Ctx:         r.Context(),
		Title:       r.PostFormValue("title"),
		User:        r.PostFormValue("user"),
		Description: r.PostFormValue("description"),
	})

	if err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(httputil.Response{
		Message: "successfully created category",
		Code:    http.StatusOK,
	})
}
