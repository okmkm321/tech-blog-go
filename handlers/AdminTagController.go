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

type TagPayload struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
	Position int    `json:"position"`
}

// all
func (app *Application) getAllTags(w http.ResponseWriter, r *http.Request) {
	tag, err := app.Models.DB.TagGetAll()
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	err = app.WriteJSON(w, http.StatusOK, tag, "tags")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}

// one
func (app *Application) getOneTag(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.Logger.Print(errors.New("invalid id parameter"))
		app.ErrorJSON(w, err)
		return
	}

	t, err := app.Models.DB.GetTag(id)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	err = app.WriteJSON(w, http.StatusOK, t, "tag")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}

// create or update
func (app *Application) editTag(w http.ResponseWriter, r *http.Request) {
	var tp TagPayload
	err := json.NewDecoder(r.Body).Decode(&tp)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}

	var tm models.Tag

	if tp.ID != 0 {
		id := tp.ID
		t, _ := app.Models.DB.GetTag(id)
		tm = *t
		tm.UpdatedAt = time.Now()
	}

	tm.ID = tp.ID
	tm.Name = tp.Name
	tm.IsPublic = tp.IsPublic
	tm.Position = tp.Position
	tm.UpdatedAt = time.Now()

	if tp.ID == 0 {
		err = app.Models.DB.TagCreate(tm)
		if err != nil {
			app.ErrorJSON(w, err)
			return
		}
	} else {
		err = app.Models.DB.TagUpdate(tm)
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
