package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/pkg/dep"
	"github.com/neutrinocorp/life-track-api/pkg/transport/handler"
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

	_ = handler.NewGetCategory(getCategory, r)

	listCategory, cleanLCat, err := dep.InjectListCategoriesQuery()
	if err != nil {
		panic(err)
	}
	defer cleanLCat()

	_ = handler.NewListCategory(listCategory, r)

	addCategory, cleanACat, err := dep.InjectAddCategoryHandler()
	if err != nil {
		panic(err)
	}
	defer cleanACat()

	_ = handler.NewAddCategory(addCategory, r)

	changeCategory, cleanCState, err := dep.InjectChangeCategoryState()
	if err != nil {
		panic(err)
	}
	defer cleanCState()

	_ = handler.NewChangeCategoryState(changeCategory, r)

	editCategory, cleanECat, err := dep.InjectEditCategory()
	if err != nil {
		panic(err)
	}
	defer cleanECat()

	_ = handler.NewEditCategory(editCategory, r)

	removeCategory, cleanRCat, err := dep.InjectRemoveCategory()
	if err != nil {
		panic(err)
	}
	defer cleanRCat()

	_ = handler.NewRemoveCategory(removeCategory, r)

	log.Print("starting http sandbox")

	panic(http.ListenAndServe(":8080", r))
}
