package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) getOneUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.Logger.Print(errors.New("invalid id paramater"))
		app.ErrorJSON(w, err)
		return
	}

	user, err := app.Models.DB.Get(id)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	err = app.WriteJSON(w, http.StatusOK, user, "user")

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

}
