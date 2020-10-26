package categoryhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/neutrinocorp/life-track-api/internal/application/category"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
)

type Edit struct {
	cmd    *category.UpdateCommandHandler
	router *mux.Router
}

// NewEdit creates an Edit handler with routing
func NewEdit(cmd *category.UpdateCommandHandler, r *mux.Router) *Edit {
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
	target, err := strconv.ParseInt(r.URL.Query().Get("target_time"), 10, 64)
	if err != nil {
		target = 0
	}

	if err = c.cmd.Invoke(category.UpdateCommand{
		Ctx:         r.Context(),
		ID:          mux.Vars(r)["id"],
		UserID:      r.URL.Query().Get("user_id"),
		Name:        r.URL.Query().Get("name"),
		Description: r.URL.Query().Get("description"),
		TargetTime:  target,
		Picture:     r.URL.Query().Get("picture"),
		State:       r.URL.Query().Get("state"),
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
