package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/users/:id", app.getOneUser)
	// カテゴリー
	router.HandlerFunc(http.MethodGet, "/categories", app.getAllCategories)
	router.HandlerFunc(http.MethodPost, "/categories", app.editCategory)
	router.HandlerFunc(http.MethodGet, "/categories/:id", app.getOneCategory)
	router.HandlerFunc(http.MethodDelete, "/categories/:id", app.deleteCategory)
	// タグ
	router.HandlerFunc(http.MethodGet, "/tags", app.getAllTags)
	router.HandlerFunc(http.MethodGet, "/tags/:id", app.getOneTag)
	router.HandlerFunc(http.MethodPost, "/tags", app.editTag)
	return router
}
