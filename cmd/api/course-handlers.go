package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type MetaData struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	PageSize    int `json:"page_size"`
	TotalCount  int `json:"total_count"`
}

func (app *application) getOneCourse(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	course, err := app.models.DB.Get(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, course, "course")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllCourses(w http.ResponseWriter, r *http.Request) {

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
	ct := r.URL.Query().Get("courseTypes")
	lngs := r.URL.Query().Get("languages")
	sjts := r.URL.Query().Get("subjects")
	insts := r.URL.Query().Get("institutions")
	ist9, _ := strconv.ParseBool(r.URL.Query().Get("isTu9"))
	isu15, _ := strconv.ParseBool(r.URL.Query().Get("isU15"))
	ha, _ := strconv.ParseBool(r.URL.Query().Get("hasArticles"))
	o := r.URL.Query().Get("orderBy")
	hla, _ := strconv.ParseBool(r.URL.Query().Get("hideLanguageArticle"))

	var cp models.CourseParams
	cp.PageNumber = pn
	cp.PageSize = ps
	cp.Languages = strings.ToLower(lngs)
	cp.Subjects = strings.ToLower(sjts)
	cp.SearchTerm = strings.ToLower(st)
	cp.CourseTypes = strings.ToLower(ct)
	cp.Institutions = strings.ToLower(insts)
	cp.IsTu9 = ist9
	cp.IsU15 = isu15
	cp.HasArticles = ha
	cp.OrderBy = strings.ToLower(o)
	cp.HideLanguageNArticle = hla

	// courses, err := app.models.DB.All(pn, ps)

	//return total count from all
	courses, count, err := app.models.DB.All(cp)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// // count have to be adjusted later depends on filters
	// count, err := app.models.DB.Count()
	// if err != nil {
	// 	app.errorJSON(w, err)
	// 	return
	// }

	var md MetaData
	md.PageSize = ps
	md.CurrentPage = pn
	md.TotalCount = count
	md.TotalPages = int(math.Ceil(float64(count) / float64(ps)))
	// app.metadata = md

	js, _ := json.Marshal(md)
	w.Header().Set("Pagination", string(js))
	w.Header().Set("Access-Control-Expose-Headers", "Pagination")

	err = app.writeJSON(w, http.StatusOK, courses, "courses")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) getFilters(w http.ResponseWriter, r *http.Request) {
	filters, err := app.models.DB.GetFilters()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, filters, "filters")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) editCourse(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(0)

	var course models.Course

	course.ID, _ = strconv.Atoi(r.FormValue("id"))
	course.UniversityId = r.FormValue("universityId")
	course.CourseType = r.FormValue("courseTypes")
	course.NameEn = r.FormValue("nameEn")
	course.NameEnShort = r.FormValue("nameEnShort")
	course.NameCh = r.FormValue("nameCh")
	course.NameChShort = r.FormValue("nameChShort")
	course.TuitionFees = r.FormValue("tuitionFees")
	course.Beginning = r.FormValue("beginning")
	course.Subject = r.FormValue("subjects")
	course.Daadlink = r.FormValue("daadlink")
	course.IsElearning, _ = strconv.ParseBool(r.FormValue("isElearning"))
	course.IsCompleteOnlinePossible, _ = strconv.ParseBool(r.FormValue("isCompleteOnlinePossible"))
	course.IsFromDaad, _ = strconv.ParseBool(r.FormValue("isFromDaad"))
	course.ProgrammeDuration = r.FormValue("programmeDuration")
	course.ApplicationDeadline = r.FormValue("applicationDeadline")
	course.CreatedAt = time.Now()

	err := app.models.DB.InsertCourse(course)
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
