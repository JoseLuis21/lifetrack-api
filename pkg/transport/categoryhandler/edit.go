package categoryhandler

import (
	"encoding/json"
	"net/http"

	"github.com/neutrinocorp/life-track-api/internal/application/category"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
)

type Edit struct {
	cmd    *category.EditHandler
	router *mux.Router
}

// NewEdit creates an Edit handler with routing
func NewEdit(cmd *category.EditHandler, r *mux.Router) *Edit {
	h := &Edit{
		cmd:    cmd,
		router: r,
	}
	h.mapRoute()

	return h
}

func (c *Edit) mapRoute() {
	c.router.Path("/category/{id}").Methods(http.MethodPut, http.MethodPatch).HandlerFunc(c.Handler)
}

func (c Edit) GetRouter() *mux.Router {
	return c.router
}

func (c Edit) Handler(w http.ResponseWriter, r *http.Request) {
	if err := c.cmd.Invoke(category.Edit{
		Ctx:         r.Context(),
		ID:          mux.Vars(r)["id"],
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
		Theme:       r.PostFormValue("theme"),
		Image:       r.PostFormValue("image"),
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
