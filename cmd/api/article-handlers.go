package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getOneArticle(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	article, err := app.models.DB.GetOneArticle(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, article, "article")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllArticles(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /v1/articles", r.URL.Query())
	pn, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ps, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	st := r.URL.Query().Get("searchTerm")
	srcs := r.URL.Query().Get("sources")
	bschs := r.URL.Query().Get("bsSchools")
	bsds := r.URL.Query().Get("bsDepartments")
	mschs := r.URL.Query().Get("msSchools")
	msds := r.URL.Query().Get("msDepartments")
	ct := r.URL.Query().Get("courseType")
	ha, _ := strconv.ParseBool(r.URL.Query().Get("hideApplication"))

	var ap models.ArticleParams
	ap.PageNumber = pn
	ap.PageSize = ps
	ap.SearchTerm = strings.ToLower(st)
	ap.Sources = strings.ToLower(srcs)
	ap.BsSchools = strings.ToLower(bschs)
	ap.BsDepartments = strings.ToLower(bsds)
	ap.MsSchools = strings.ToLower(mschs)
	ap.MsDepartments = strings.ToLower(msds)
	ap.CourseType = ct
	ap.HideApplication = ha

	articles, count, err := app.models.DB.GetArticles(ap)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var md MetaData
	md.PageSize = ps
	md.CurrentPage = pn
	md.TotalCount = count
	md.TotalPages = int(math.Ceil(float64(count) / float64(ps)))

	js, _ := json.Marshal(md)
	w.Header().Set("Pagination", string(js))
	w.Header().Set("Access-Control-Expose-Headers", "Pagination")

	err = app.writeJSON(w, http.StatusOK, articles, "articles")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getArticleFilters(w http.ResponseWriter, r *http.Request) {
	filters, err := app.models.DB.GetArticleFilters()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, filters, "articleFilters")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) editArticle(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(0)

	var ca models.CourseArticle

	ca.ArticleID, _ = strconv.Atoi(r.FormValue("id"))
	ca.CourseID, _ = strconv.Atoi(r.FormValue("courseId"))
	ca.Result = r.FormValue("result")
	ca.IsDecision, _ = strconv.ParseBool(r.FormValue("isDecision"))

	err := app.models.DB.InsertArticle(ca)
	if err != nil {
		app.errorJSON(w, err)
		return
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
