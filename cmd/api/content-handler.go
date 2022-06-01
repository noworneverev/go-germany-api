package main

import (
	"backend/models"
	"net/http"
	"strconv"
	"time"
)

func (app *application) editContent(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	var c models.Content

	layout := "2006-01-02"

	c.ID, _ = strconv.Atoi(r.FormValue("id"))
	c.Link = r.FormValue("link")
	c.Title = r.FormValue("title")
	c.Author = r.FormValue("author")
	c.PublishedAt, _ = time.Parse(layout, r.FormValue("publishedAt"))
	// t, _ := time.Parse(layout, r.FormValue("publishedAd"))
	// c.PublishedAt = t.Format("2006-01-02")
	c.Source = r.FormValue("source")
	c.AuthorBsSchool = r.FormValue("authorBsSchool")
	c.AuthorBsSchoolShort = r.FormValue("authorBsSchoolShort")
	c.AuthorBsDepartment = r.FormValue("authorBsDepartment")
	c.AuthorBsGpa = r.FormValue("authorBsGpa")
	c.AuthorMsSchool = r.FormValue("authorMsSchool")
	c.AuthorMsSchoolShort = r.FormValue("authorMsSchoolShort")
	c.AuthorMsDepartment = r.FormValue("authorMsDepartment")
	c.AuthorMsGpa = r.FormValue("authorMsGpa")
	c.AuthorToefl = r.FormValue("authorToefl")
	c.AuthorIelts = r.FormValue("authorIelts")
	c.AuthorGre = r.FormValue("authorGre")
	c.AuthorGmat = r.FormValue("authorGmat")
	c.AuthorTestdaf = r.FormValue("authorTestdaf")
	c.AuthorGoethe = r.FormValue("authorGoethe")
	c.CourseType = r.FormValue("courseType")

	err := app.models.DB.InsertContent(c)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// type jsonResp struct {
	// 	OK bool `json:"ok"`
	// }

	ok := jsonResp{
		OK: true,
	}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
