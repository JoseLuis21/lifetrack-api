package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/pkg/dep"
	"github.com/neutrinocorp/life-track-api/pkg/transport/categoryhandler"
)

func main() {
	r := mux.NewRouter()
	r = r.PathPrefix("/live").Subrouter()
	// Add middlewares

	// Known code-smell, ignore
	getCategory, cleanCGet, err := dep.InjectGetCategoryQuery()
	if err != nil {
		panic(err)
	}
	defer cleanCGet()

	_ = categoryhandler.NewGet(getCategory, r)

	listCategory, cleanLCat, err := dep.InjectListCategoriesQuery()
	if err != nil {
		panic(err)
	}
	defer cleanLCat()

	_ = categoryhandler.NewList(listCategory, r)

	addCategory, cleanACat, err := dep.InjectAddCategoryHandler()
	if err != nil {
		panic(err)
	}
	defer cleanACat()

	_ = categoryhandler.NewAdd(addCategory, r)

	editCategory, cleanECat, err := dep.InjectEditCategory()
	if err != nil {
		panic(err)
	}
	defer cleanECat()

	_ = categoryhandler.NewEdit(editCategory, r)

	removeCategory, cleanRCat, err := dep.InjectRemoveCategory()
	if err != nil {
		panic(err)
	}
	defer cleanRCat()

	_ = categoryhandler.NewRemove(removeCategory, r)

	log.Print("starting http sandbox")

	panic(http.ListenAndServe(":8080", r))
}
