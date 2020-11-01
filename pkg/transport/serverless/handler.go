package serverless

import "github.com/gorilla/mux"

// Handler serverless handler which contains a mux router
type Handler interface {
	GetRouter() *mux.Router
}
