package categoryhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/neutrinocorp/life-track-api/internal/application/category"

	"github.com/alexandria-oss/common-go/exception"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
)

type ChangeState struct {
	cmd    *category.ChangeStateHandler
	router *mux.Router
}

// NewChangeState creates a ChangeState handler with routing
func NewChangeState(cmd *category.ChangeStateHandler, r *mux.Router) *ChangeState {
	h := &ChangeState{
		cmd:    cmd,
		router: r,
	}
	h.mapRoute()

	return h
}

func (c *ChangeState) mapRoute() {
	c.router.Path("/category/{id}/state").Methods(http.MethodPatch, http.MethodPut).HandlerFunc(c.Handler)
}

func (c ChangeState) GetRouter() *mux.Router {
	return c.router
}

func (c ChangeState) Handler(w http.ResponseWriter, r *http.Request) {
	state, err := strconv.ParseBool(r.PostFormValue("state"))
	if err != nil {
		httputil.RespondErrorJSON(exception.NewFieldFormat("state", "boolean"), w)
		return
	}

	if err = c.cmd.Invoke(category.ChangeState{
		Ctx:   r.Context(),
		ID:    mux.Vars(r)["id"],
		State: state,
	}); err != nil {
		httputil.RespondErrorJSON(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(httputil.Response{
		Message: "successfully changed category state",
		Code:    http.StatusOK,
	})
}
