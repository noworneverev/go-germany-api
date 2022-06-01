package main

import (
	"backend/models"
	"net/http"
	"strconv"
	"time"
)

func (app *application) editUniversity(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(0)

	var university models.University

	university.ID, _ = strconv.Atoi(r.FormValue("id"))
	university.NameEn = r.FormValue("nameEn")
	university.NameCh = r.FormValue("nameCh")
	university.City = r.FormValue("city")
	university.IsFromDaad, _ = strconv.ParseBool(r.FormValue("isFromDaad"))
	university.IsTu9, _ = strconv.ParseBool(r.FormValue("isTu9"))
	university.IsU15, _ = strconv.ParseBool(r.FormValue("isU15"))
	university.CreatedAt = time.Now()
	university.Link = r.FormValue("link")
	university.QsRanking, _ = strconv.Atoi(r.FormValue("qsRanking"))

	err := app.models.DB.InsertUniversity(university)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	type jsonResp struct {
		OK bool `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
