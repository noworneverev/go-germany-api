package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// GetOneArticle returns one course and error, if any
func (m *DBModel) GetOneArticle(id int) (*Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	query := `select c.id, c.link, c.title, c.author, c.published_date, c.source, 
	c.author_bs_school, c.author_bs_school_short, c.author_bs_department, c.author_bs_gpa,
	c.author_ms_school, c.author_ms_school_short, c.author_ms_department, c.author_ms_gpa,
	c.author_toefl, c.author_ielts, c.author_gre, c.author_gmat, c.author_testdaf, c.author_goethe, c.course_type
	from content c
	where c.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	article, err := ScanArticle(row)
	if err != nil {
		return nil, err
	}

	// get the courses
	err = article.SetCourses(m, ctx)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func ScanArticle(row *sql.Row) (Article, error) {
	var article Article
	err := row.Scan(
		&article.ID,
		&article.Link,
		&article.Title,
		&article.Author,
		&article.PublishedAt,
		&article.Source,
		&article.AuthorBsSchool,
		&article.AuthorBsSchoolShort,
		&article.AuthorBsDepartment,
		&article.AuthorBsGpa,
		&article.AuthorMsSchool,
		&article.AuthorMsSchoolShort,
		&article.AuthorMsDepartment,
		&article.AuthorMsGpa,
		&article.AuthorToefl,
		&article.AuthorIelts,
		&article.AuthorGre,
		&article.AuthorGmat,
		&article.AuthorTestdaf,
		&article.AuthorGoethe,
		&article.CourseType,
	)
	return article, err
}

func (article *Article) SetCourses(m *DBModel, ctx context.Context) error {

	courseQuery := `select a.course_id, a.result, a.is_decision, c.name_en, 
	c.daadlink, c.is_from_daad, c.course_type, c.programme_duration,
	c.tuition_fees, c.beginning, c.subject, c.application_deadline,
	u.name_en, u.name_ch, u.link, u.is_tu9, u.is_u15, u.city
	from article as a
	left join course as c on c.id = a.course_id 
	left join university as u on u.id = c.university_id
	where a.id = $1`

	courseRows, _ := m.DB.QueryContext(ctx, courseQuery, article.ID)

	// var courses []Course
	var articleCourses []ArticleCourse

	for courseRows.Next() {
		var ac ArticleCourse
		err := courseRows.Scan(
			&ac.Course.ID,
			&ac.Result,
			&ac.IsDecision,
			&ac.Course.NameEn,
			&ac.Course.Daadlink,
			&ac.Course.IsFromDaad,
			&ac.Course.CourseType,
			&ac.Course.ProgrammeDuration,
			&ac.Course.TuitionFees,
			&ac.Course.Beginning,
			&ac.Course.Subject,
			&ac.Course.ApplicationDeadline,
			&ac.Course.UniversityNameEn,
			&ac.Course.UniversityNameCh,
			&ac.Course.UniversityLink,
			&ac.Course.IsTu9,
			&ac.Course.IsU15,
			&ac.Course.City,
		)
		if err != nil {
			return err
		}

		// courses = append(courses, ac.Course)
		articleCourses = append(articleCourses, ac)
	}
	courseRows.Close()

	article.ArticleCourse = articleCourses
	return nil
}

func (m *DBModel) GetArticles(ap ArticleParams) ([]*Article, int, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var rows *sql.Rows
	var err error

	var query string
	count := 0

	baseQueryString := `select c.id, c.link, c.title, c.author, c.published_date, c.source, 
	c.author_bs_school, c.author_bs_school_short, c.author_bs_department, c.author_bs_gpa,
	c.author_ms_school, c.author_ms_school_short, c.author_ms_department, c.author_ms_gpa,
	c.author_toefl, c.author_ielts, c.author_gre, c.author_gmat, c.author_testdaf, c.author_goethe, c.course_type
	from content c`

	srcarr := strings.Split(ap.Sources, ",")
	bsharr := strings.Split(ap.BsSchools, ",")
	bsdarr := strings.Split(ap.BsDepartments, ",")
	msharr := strings.Split(ap.MsSchools, ",")
	msdarr := strings.Split(ap.MsDepartments, ",")

	srcs := "'" + strings.Join(srcarr, "','") + "'"
	bshs := "'" + strings.Join(bsharr, "','") + "'"
	bsds := "'" + strings.Join(bsdarr, "','") + "'"
	mshs := "'" + strings.Join(msharr, "','") + "'"
	msds := "'" + strings.Join(msdarr, "','") + "'"

	likeSearchTerm := ""
	inSrcs := ""
	inBsSchools := ""
	inBsDepartments := ""
	inMsSchools := ""
	inMsDepartments := ""
	equalCourseType := ""

	orderBy := "order by c.published_date desc"
	limitNOffset := fmt.Sprintf("limit %d offset %d", ap.PageSize, (ap.PageNumber-1)*ap.PageSize)

	var whereArr []string
	where := ""

	if len(ap.SearchTerm) > 0 {
		likeSearchTerm = fmt.Sprintf(`(lower(c.author) like '%%%[1]s%%' or lower(c.title) like '%%%[1]s%%' or lower(c.source) like '%%%[1]s%%' or 
		lower(c.author_bs_school) like '%%%[1]s%%' or lower(c.author_bs_school_short) like '%%%[1]s%%' or lower(c.author_bs_department) like '%%%[1]s%%' or 
		lower(c.author_ms_school) like '%%%[1]s%%' or lower(c.author_ms_school_short) like '%%%[1]s%%' or lower(c.author_ms_department) like '%%%[1]s%%')`, ap.SearchTerm)
		whereArr = append(whereArr, likeSearchTerm)
	}

	if len(ap.Sources) > 0 {
		inSrcs = "(lower(c.source) in " + "(" + srcs + "))"
		whereArr = append(whereArr, inSrcs)
	}

	if len(ap.BsSchools) > 0 {
		inBsSchools = "(lower(c.author_bs_school_short) in " + "(" + bshs + "))"
		whereArr = append(whereArr, inBsSchools)
	}

	if len(ap.BsDepartments) > 0 {
		inBsDepartments = "(lower(c.author_bs_department) in " + "(" + bsds + "))"
		whereArr = append(whereArr, inBsDepartments)
	}

	if len(ap.MsSchools) > 0 {
		inMsSchools = "(lower(c.author_ms_school_short) in " + "(" + mshs + "))"
		whereArr = append(whereArr, inMsSchools)
	}

	if len(ap.MsDepartments) > 0 {
		inMsDepartments = "(lower(c.author_ms_department) in " + "(" + msds + "))"
		whereArr = append(whereArr, inMsDepartments)
	}

	if len(ap.CourseType) > 0 {
		equalCourseType = fmt.Sprintf("(c.course_type = %s )", ap.CourseType)
		whereArr = append(whereArr, equalCourseType)
	}

	where = strings.Join(whereArr, " and ")

	if len(whereArr) > 0 {
		query = fmt.Sprintf("%s where %s", baseQueryString, where)
	} else {
		query = baseQueryString
	}

	countQuery := fmt.Sprintf("select count(*) from (%s) as c", query)

	row := m.DB.QueryRowContext(ctx, countQuery)
	err = row.Scan(&count)
	if err != nil {
		return nil, -1, err
	}

	//query with limit and offset
	query = fmt.Sprintf("%s %s %s", query, orderBy, limitNOffset)

	rows, err = m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, -1, err
	}
	defer rows.Close()

	var articles []*Article
	for rows.Next() {
		article, err := ScanArticles(rows)
		if err != nil {
			return nil, -1, err
		}

		if !ap.HideApplication {
			// get the courses
			err = article.SetCourses(m, ctx)
			if err != nil {
				return nil, -1, err
			}
		}

		articles = append(articles, &article)
	}

	return articles, count, nil
}

func ScanArticles(rows *sql.Rows) (Article, error) {
	var article Article
	err := rows.Scan(
		&article.ID,
		&article.Link,
		&article.Title,
		&article.Author,
		&article.PublishedAt,
		&article.Source,
		&article.AuthorBsSchool,
		&article.AuthorBsSchoolShort,
		&article.AuthorBsDepartment,
		&article.AuthorBsGpa,
		&article.AuthorMsSchool,
		&article.AuthorMsSchoolShort,
		&article.AuthorMsDepartment,
		&article.AuthorMsGpa,
		&article.AuthorToefl,
		&article.AuthorIelts,
		&article.AuthorGre,
		&article.AuthorGmat,
		&article.AuthorTestdaf,
		&article.AuthorGoethe,
		&article.CourseType,
	)
	return article, err
}

func (m *DBModel) GetArticleFilters() (*ArticleFilters, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var articleFilters ArticleFilters

	var sources []string
	var bsSchools []string
	var bsDepartments []string
	var msSchools []string
	var msDepartments []string
	var courseTypes []string

	// sources
	query := `select source from content
	group by source
	order by case when source = 'PTT' then 1
					when source = 'FB' then 2
					when source = 'Dcard' then 3
					when source = 'Medium' then 4
					when source = 'Blog' then 5
					else 6
					end asc`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var src string
		err := rows.Scan(
			&src,
		)

		if err != nil {
			return nil, err
		}
		sources = append(sources, src)
	}
	articleFilters.Sources = sources

	//Bs schools
	query = "select distinct author_bs_school_short from content where author_bs_school_short <> '' order by author_bs_school_short"
	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var bsl string
		err := rows.Scan(
			&bsl,
		)

		if err != nil {
			return nil, err
		}

		bsSchools = append(bsSchools, bsl)
	}

	articleFilters.BsSchools = bsSchools

	//Bs Departments
	query = "select distinct author_bs_department from content where author_bs_department <> '' order by author_bs_department"
	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var bsd string
		err := rows.Scan(
			&bsd,
		)

		if err != nil {
			return nil, err
		}

		bsDepartments = append(bsDepartments, bsd)
	}
	articleFilters.BsDepartments = bsDepartments

	//Ms schools
	query = "select distinct author_ms_school_short from content where author_ms_school_short <> '' order by author_ms_school_short"
	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var msl string
		err := rows.Scan(
			&msl,
		)

		if err != nil {
			return nil, err
		}

		msSchools = append(msSchools, msl)
	}
	articleFilters.MsSchools = msSchools

	//Ms Departments
	query = "select distinct author_ms_department from content where author_ms_department <> '' order by author_ms_department"
	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var msd string
		err := rows.Scan(
			&msd,
		)

		if err != nil {
			return nil, err
		}

		msDepartments = append(msDepartments, msd)
	}
	articleFilters.MsDepartments = msDepartments

	//Course types
	query = "select distinct course_type from content"
	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ct string
		err := rows.Scan(
			&ct,
		)

		if err != nil {
			return nil, err
		}

		courseTypes = append(courseTypes, ct)
	}
	articleFilters.CourseTypes = courseTypes

	return &articleFilters, nil
}

func (m *DBModel) InsertArticle(ca CourseArticle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	stmt := `insert into article (id, course_id, result, is_decision) values ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, stmt,
		ca.ArticleID,
		ca.CourseID,
		ca.Result,
		ca.IsDecision,
	)
	if err != nil {
		return err
	}
	return nil
}
