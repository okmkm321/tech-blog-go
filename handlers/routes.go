package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/api/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/api/users/:id", app.getOneUser)
	// カテゴリー
	router.HandlerFunc(http.MethodGet, "/api/categories", app.getAllCategories)
	router.HandlerFunc(http.MethodPost, "/api/categories", app.editCategory)
	router.HandlerFunc(http.MethodGet, "/api/categories/:id", app.getOneCategory)
	router.HandlerFunc(http.MethodDelete, "/api/categories/:id", app.deleteCategory)
	// タグ
	router.HandlerFunc(http.MethodGet, "/api/tags", app.getAllTags)
	router.HandlerFunc(http.MethodPost, "/api/tags", app.editTag)
	router.HandlerFunc(http.MethodGet, "/api/tags/:id", app.getOneTag)
	router.HandlerFunc(http.MethodDelete, "/api/tags/:id", app.deleteTag)
	// ブログ
	router.HandlerFunc(http.MethodGet, "/api/blogs", app.getAllBlogs)
	router.HandlerFunc(http.MethodPost, "/api/blogs", app.editBlog)
	router.HandlerFunc(http.MethodGet, "/api/blogs/:id", app.getOneBlog)
	router.HandlerFunc(http.MethodPost, "/api/blogs/eyecatch", app.eyeCatchUpload)
	return router
}
