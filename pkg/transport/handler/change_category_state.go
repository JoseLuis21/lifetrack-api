package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexandria-oss/common-go/exception"

	"github.com/alexandria-oss/common-go/httputil"
	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
)

type ChangeCategoryState struct {
	cmd    *command.ChangeCategoryStateHandler
	router *mux.Router
}

// NewChangeCategoryState creates an change category state handler with routing
func NewChangeCategoryState(cmd *command.ChangeCategoryStateHandler, r *mux.Router) *ChangeCategoryState {
	h := &ChangeCategoryState{
		cmd:    cmd,
		router: r,
	}
	h.mapRoute()

	return h
}

func (c *ChangeCategoryState) mapRoute() {
	c.router.Path("/category/{id}/state").Methods(http.MethodPatch, http.MethodPut).HandlerFunc(c.Handler)
}

func (c ChangeCategoryState) GetRouter() *mux.Router {
	return c.router
}

func (c ChangeCategoryState) Handler(w http.ResponseWriter, r *http.Request) {
	state, err := strconv.ParseBool(r.PostFormValue("state"))
	if err != nil {
		httputil.RespondErrorJSON(exception.NewFieldFormat("state", "boolean"), w)
		return
	}

	if err = c.cmd.Invoke(command.ChangeCategoryState{
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
