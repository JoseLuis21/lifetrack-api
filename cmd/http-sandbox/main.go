package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neutrinocorp/life-track-api/pkg/dep"
	"github.com/neutrinocorp/life-track-api/pkg/transport/handler"
)

func main() {
	r := mux.NewRouter()

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

	panic(http.ListenAndServe(":8080", r))
}
