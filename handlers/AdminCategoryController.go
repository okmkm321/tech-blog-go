package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"project/tech-blog-go/models"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type CategoryPayload struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	IsPublic bool        `json:"is_public"`
	ParentId interface{} `json:"parent_id"`
}

// all
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

// one
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

// create or update
func (app *Application) editCategory(w http.ResponseWriter, r *http.Request) {
	var cp CategoryPayload
	err := json.NewDecoder(r.Body).Decode(&cp)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}

	var cm models.Category

	if cp.ID != 0 {
		id := cp.ID
		c, _ := app.Models.DB.GetCategory(id)
		cm = *c
		cm.UpdatedAt = time.Now()
	}

	cm.ID = cp.ID
	cm.Name = cp.Name
	cm.Slug = cp.Slug
	cm.IsPublic = cp.IsPublic
	cm.ParentId = cp.ParentId
	cm.UpdatedAt = time.Now()

	if cp.ID == 0 {
		err = app.Models.DB.CategoryCreate(cm)
		if err != nil {
			app.ErrorJSON(w, err)
			return
		}
	} else {
		err = app.Models.DB.CategoryUpdate(cm)
		if err != nil {
			app.ErrorJSON(w, err)
			return
		}
	}

	ok := JsonResp{
		OK: true,
	}

	err = app.WriteJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}

func (app *Application) deleteCategory(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	err = app.Models.DB.CategoryDelete(id)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	ok := JsonResp{
		OK: true,
	}

	err = app.WriteJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}
