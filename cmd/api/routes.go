package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodPost, "/v1/account/signin", app.signIn)

	router.HandlerFunc(http.MethodGet, "/v1/course/:id", app.getOneCourse)
	router.HandlerFunc(http.MethodGet, "/v1/courses", app.getAllCourses)
	router.HandlerFunc(http.MethodGet, "/v1/courses/filters", app.getFilters)

	router.HandlerFunc(http.MethodGet, "/v1/article/:id", app.getOneArticle)
	router.HandlerFunc(http.MethodGet, "/v1/articles", app.getAllArticles)
	router.HandlerFunc(http.MethodGet, "/v1/articles/filters", app.getArticleFilters)

	router.POST("/v1/admin/editcourse", app.wrap(secure.ThenFunc(app.editCourse)))
	// router.HandlerFunc(http.MethodPost, "/v1/admin/editcourse", app.editCourse)

	router.POST("/v1/admin/edituniversity", app.wrap(secure.ThenFunc(app.editUniversity)))
	router.POST("/v1/admin/editcontent", app.wrap(secure.ThenFunc(app.editContent)))
	router.POST("/v1/admin/editarticle", app.wrap(secure.ThenFunc(app.editArticle)))

	return app.enableCORS(router)
}
