package categoryhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/category"
)

type Add struct {
	cmd    *category.AddCommandHandler
	router *mux.Router
}

// NewAdd creates an Add handler with routing
func NewAdd(cmd *category.AddCommandHandler, r *mux.Router) *Add {
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
	id, err := c.cmd.Invoke(category.AddCommand{
		Ctx:    r.Context(),
		UserID: r.PostFormValue("user_id"),
		Name:   r.PostFormValue("name"),
	})
	if err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(httputil.Response{
		Message: fmt.Sprintf("successfully created category %s", id),
		Code:    http.StatusOK,
	})
}
