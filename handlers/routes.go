package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/users/:id", app.getOneUser)
	router.HandlerFunc(http.MethodGet, "/categories", app.getAllCategories)
	router.HandlerFunc(http.MethodGet, "/categories/:id", app.getOneCategory)
	router.HandlerFunc(http.MethodPost, "/categories", app.insertCategory)
	return router
}
