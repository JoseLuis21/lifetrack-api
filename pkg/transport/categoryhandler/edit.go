package categoryhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/neutrinocorp/lifetrack-api/internal/application/category"

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
	target, err := strconv.ParseInt(r.PostFormValue("target_time"), 10, 64)
	if err != nil {
		target = 0
	}

	if err = c.cmd.Invoke(category.UpdateCommand{
		Ctx:         r.Context(),
		ID:          mux.Vars(r)["id"],
		UserID:      r.PostFormValue("user_id"),
		Name:        r.PostFormValue("name"),
		Description: r.PostFormValue("description"),
		TargetTime:  target,
		Picture:     r.PostFormValue("picture"),
		State:       r.PostFormValue("state"),
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
