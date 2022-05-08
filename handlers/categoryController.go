package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"project/tech-blog-go/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type CategoryPayload struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	State    int    `json:"state"`
	ParentId int    `json:"parent_id"`
}

func (app *Application) getAllCategories(w http.ResponseWriter, r *http.Request) {
	ctg, err := app.Models.DB.CategoryGetAll()
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, ctg, "categories")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}

func (app *Application) getOneCategory(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.Logger.Print(errors.New("invalid id parameter"))
		app.ErrorJSON(w, err)
		return
	}

	ctg, err := app.Models.DB.GetCategory(id)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	err = app.WriteJSON(w, http.StatusOK, ctg, "category")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}

func (app *Application) insertCategory(w http.ResponseWriter, r *http.Request) {
	var cp CategoryPayload
	err := json.NewDecoder(r.Body).Decode(&cp)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}

	var cm models.Category
	cm.Name = cp.Name
	cm.Slug = cp.Slug
	cm.State = cp.State
	cm.ParentId = cp.ParentId
	// cm.State, _ = strconv.Atoi(cp.State)
	// cm.ParentId, _ = strconv.Atoi(cp.ParentId)

	err = app.Models.DB.CategoryCreate(cm)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	ok := jsonResp{
		OK: true,
	}

	err = app.WriteJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}
