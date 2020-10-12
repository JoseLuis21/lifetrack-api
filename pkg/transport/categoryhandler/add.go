package categoryhandler

import (
	"encoding/json"
	"net/http"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/category"
)

type Add struct {
	cmd    *category.AddHandler
	router *mux.Router
}

// NewAdd creates an Add handler with routing
func NewAdd(cmd *category.AddHandler, r *mux.Router) *Add {
	h := &Add{
		cmd:    cmd,
		router: r,
	}
	h.mapRoute()

	return h
}

func (c *Add) mapRoute() {
	c.router.Path("/category").Methods(http.MethodPost).HandlerFunc(c.Handler)
}

func (c Add) GetRouter() *mux.Router {
	return c.router
}

func (c Add) Handler(w http.ResponseWriter, r *http.Request) {
	if err := c.cmd.Invoke(category.Add{
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
