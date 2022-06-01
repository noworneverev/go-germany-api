package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get returns one course and error, if any
func (m *DBModel) Get(id int) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	query := `select c.id, c.university_id, c.course_type, c.name_en, c.name_en_short, c.tuition_fees, c.beginning, c.subject, c.daadlink, c.is_elearning, c.application_deadline,
	c.is_complete_online_possible, c.programme_duration, c.is_from_daad, c.created_at, COALESCE(c.updated_at, c.created_at),
	u.name_en, u.name_ch, u.city, u.is_tu9, u.is_u15, COALESCE(u.qs_ranking, 0), u.link
	from course as c
	left join university as u on c.university_id = u.id
	where c.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	course, err := ScanCourse(row)
	if err != nil {
		return nil, err
	}

	// get the languages
	err = course.SetLanguages(m, ctx)
	if err != nil {
		return nil, err
	}

	// get the articles
	err = course.SetArticles(m, ctx)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

// Count return length of courses
func (m *DBModel) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	query := `select count(*) from course`
	row := m.DB.QueryRowContext(ctx, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

// Filters return filter object
func (m *DBModel) GetFilters() (*Filters, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var filters Filters

	var courseTypes []string
	var languages []string
	var subjects []string
	var institutions []string

	// course types
	query := "select distinct course_type from course order by course_type"
	rows, err := m.DB.QueryContext(ctx, query)
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
	filters.CourseTypes = courseTypes

	// languages
	query = "select name from language"
	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var l string
		err := rows.Scan(
			&l,
		)

		if err != nil {
			return nil, err
		}

		languages = append(languages, l)
	}

	filters.Languages = languages

	// subjects
	query = "select distinct subject from course where subject <> '' order by subject"
	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var s string
		err := rows.Scan(
			&s,
		)

		if err != nil {
			return nil, err
		}

		subjects = append(subjects, s)
	}

	filters.Subjects = subjects

	// institutions
	query = "select distinct name_en from university order by name_en"
	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var i string
		err := rows.Scan(
			&i,
		)

		if err != nil {
			return nil, err
		}

		institutions = append(institutions, i)
	}

	filters.Institutions = institutions

	return &filters, nil
}

// All return all courses and error, if any
// func (m *DBModel) All(pageNumber int, pageSize int) ([]*Course, error) {
func (m *DBModel) All(cp CourseParams) ([]*Course, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var rows *sql.Rows
	var err error

	var query string
	count := 0
	// baseQueryString := `select c.id, c.university_id, c.course_type, c.name_en, c.name_en_short, c.tuition_fees, c.beginning, c.subject, c.daadlink, c.is_elearning, c.application_deadline,
	// c.is_complete_online_possible, c.programme_duration, c.is_from_daad, c.created_at, COALESCE(c.updated_at, c.created_at),
	// u.name_en, u.name_ch, u.city, u.is_tu9, u.is_u15, COALESCE(u.qs_ranking, 0), u.link
	// from course as c
	// left join university as u on c.university_id = u.id`

	baseQueryString := `select c.id, c.university_id, c.course_type, c.name_en, c.name_en_short, c.tuition_fees, c.beginning, c.subject, c.daadlink, c.is_elearning, c.application_deadline,
	c.is_complete_online_possible, c.programme_duration, c.is_from_daad, c.created_at, COALESCE(c.updated_at, c.created_at),
	u.name_en, u.name_ch, u.city, u.is_tu9, u.is_u15, COALESCE(u.qs_ranking, 0), u.link, COALESCE(string_agg(distinct  cl.name, ','),'') as languages, count(a.course_id) as article_count
	from course as c
	left join university as u on c.university_id = u.id	 
	left join (select cl.course_id, l.name
				from courses_languages as cl
				left join language as l on cl.language_id = l.id ) as cl on c.id = cl.course_id
	left join article as a on a.course_id = c.id`

	csarr := strings.Split(cp.CourseTypes, ",")
	usarr := strings.Split(cp.Institutions, ",")
	// subarr := strings.Split(cp.Subjects, ",")
	subarr := strings.Split(cp.Subjects, ";")
	lngarr := strings.Split(cp.Languages, ",")

	// log.Println(cp.Subjects)

	cs := "'" + strings.Join(csarr, "','") + "'"
	us := "'" + strings.Join(usarr, "','") + "'"
	subs := "'" + strings.Join(subarr, "','") + "'"

	inCourse := ""
	likeSearchTerm := ""
	inUniversity := ""
	inSubject := ""
	havingLngs := ""
	isTu9 := ""
	isU15 := ""
	hasArticles := ""

	groupBy := `group by c.id, c.university_id, c.course_type, c.name_en, c.name_en_short, c.tuition_fees, c.beginning, c.subject, c.daadlink, c.is_elearning, c.application_deadline,
	c.is_complete_online_possible, c.programme_duration, c.is_from_daad, c.created_at, COALESCE(c.updated_at, c.created_at),
	u.name_en, u.name_ch, u.city, u.is_tu9, u.is_u15, COALESCE(u.qs_ranking, 0), u.link, a.course_id`
	orderBy := "order by u.name_en, c.course_type, c.name_en"
	limitNOffset := fmt.Sprintf("limit %d offset %d", cp.PageSize, (cp.PageNumber-1)*cp.PageSize)

	// where likeSearchTerm inCourse inUniversity inSubject isTu9 isU15
	// having havingLngs hasArticles
	var whereArr []string
	var havingArr []string
	where := ""
	having := ""

	if len(cp.SearchTerm) > 0 {
		likeSearchTerm = fmt.Sprintf(`(lower(c.name_en) like '%%%[1]s%%' or lower(u.name_en) like '%%%[1]s%%' or lower(u.name_ch) like '%%%[1]s%%' or lower(c.subject) like '%%%[1]s%%')`, cp.SearchTerm)
		whereArr = append(whereArr, likeSearchTerm)
	}

	if len(cp.CourseTypes) > 0 {
		inCourse = "(c.course_type in " + "(" + cs + "))"
		whereArr = append(whereArr, inCourse)
	}

	if len(cp.Institutions) > 0 {
		inUniversity = "(lower(u.name_en) in " + "(" + us + "))"
		whereArr = append(whereArr, inUniversity)
	}

	if len(cp.Subjects) > 0 {
		inSubject = "(lower(c.subject) in " + "(" + subs + "))"
		whereArr = append(whereArr, inSubject)
	}

	if cp.IsTu9 {
		isTu9 = "(u.is_tu9)"
		whereArr = append(whereArr, isTu9)
	}

	if cp.IsU15 {
		isU15 = "(u.is_u15)"
		whereArr = append(whereArr, isU15)
	}

	if len(cp.Languages) > 0 {
		for i, lng := range lngarr {
			if i == len(lngarr)-1 {
				havingLngs += "'" + lng + "'" + "= ANY(string_to_array(lower(string_agg(distinct cl.name, ',')), ','))"
			} else {
				havingLngs += "'" + lng + "'" + "= ANY(string_to_array(lower(string_agg(distinct cl.name, ',')), ',')) or "
			}
		}
		havingLngs = "(" + havingLngs + ")"
		havingArr = append(havingArr, havingLngs)
	}

	if cp.HasArticles {
		hasArticles = "count(a.course_id) > 0"
		havingArr = append(havingArr, hasArticles)
	}

	where = strings.Join(whereArr, " and ")
	having = strings.Join(havingArr, " and ")

	if len(whereArr) > 0 && len(havingArr) > 0 {
		query = fmt.Sprintf("%s where %s %s having %s", baseQueryString, where, groupBy, having)
	} else if len(whereArr) > 0 {
		query = fmt.Sprintf("%s where %s %s", baseQueryString, where, groupBy)
	} else if len(havingArr) > 0 {
		query = fmt.Sprintf("%s %s having %s", baseQueryString, groupBy, having)
	} else {
		query = fmt.Sprintf("%s %s", baseQueryString, groupBy)
	}

	//original query to count total rows
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

	var courses []*Course
	for rows.Next() {
		course, err := ScanCourses(rows)
		if err != nil {
			return nil, -1, err
		}

		if !cp.HideLanguageNArticle {
			// get the languages
			err = course.SetLanguages(m, ctx)
			if err != nil {
				return nil, -1, err
			}

			// get the articles
			err = course.SetArticles(m, ctx)
			if err != nil {
				return nil, -1, err
			}
		}

		courses = append(courses, &course)
	}
	// log.Println(count)

	return courses, count, nil

}

func ScanCourse(row *sql.Row) (Course, error) {
	var course Course
	err := row.Scan(
		&course.ID,
		&course.UniversityId,
		&course.CourseType,
		&course.NameEn,
		&course.NameEnShort,
		&course.TuitionFees,
		&course.Beginning,
		&course.Subject,
		&course.Daadlink,
		&course.IsElearning,
		&course.ApplicationDeadline,
		&course.IsCompleteOnlinePossible,
		&course.ProgrammeDuration,
		&course.IsFromDaad,
		&course.CreatedAt,
		&course.UpdatedAt,
		&course.UniversityNameEn,
		&course.UniversityNameCh,
		&course.City,
		&course.IsTu9,
		&course.IsU15,
		&course.QsRanking,
		&course.UniversityLink,
	)
	return course, err
}

func ScanCourses(row *sql.Rows) (Course, error) {
	var course Course
	err := row.Scan(
		&course.ID,
		&course.UniversityId,
		&course.CourseType,
		&course.NameEn,
		&course.NameEnShort,
		&course.TuitionFees,
		&course.Beginning,
		&course.Subject,
		&course.Daadlink,
		&course.IsElearning,
		&course.ApplicationDeadline,
		&course.IsCompleteOnlinePossible,
		&course.ProgrammeDuration,
		&course.IsFromDaad,
		&course.CreatedAt,
		&course.UpdatedAt,
		&course.UniversityNameEn,
		&course.UniversityNameCh,
		&course.City,
		&course.IsTu9,
		&course.IsU15,
		&course.QsRanking,
		&course.UniversityLink,
		&course.Languages,
		&course.ArticleCount,
	)
	return course, err
}

func (course *Course) SetLanguages(m *DBModel, ctx context.Context) error {
	// get the languages
	languageQuery := `select
		cl.id, cl.course_id, cl.language_id, l.name
		from courses_languages as cl
		left join language as l on (l.id = cl.language_id)
		where cl.course_id = $1`
	languageRows, _ := m.DB.QueryContext(ctx, languageQuery, course.ID)
	// languages := make(map[int]string)
	var languages []string
	for languageRows.Next() {
		var cl CourseLanguage
		err := languageRows.Scan(
			&cl.ID,
			&cl.CourseID,
			&cl.LanguageID,
			&cl.Language.LanguageName,
		)
		if err != nil {
			return err
		}
		// languages[cl.ID] = cl.Language.LanguageName
		languages = append(languages, cl.Language.LanguageName)
	}
	languageRows.Close()

	course.CourseLanguage = languages
	return nil
}

func (course *Course) SetArticles(m *DBModel, ctx context.Context) error {

	articleQuery := `select
		a.id, a.course_id, a.result, a.is_decision, ct.title, ct.author, ct.link, ct.published_date, ct.source,
		ct.author_bs_school, ct.author_bs_school_short, ct.author_bs_department, ct.author_bs_gpa,
		ct.author_ms_school, ct.author_ms_school_short, ct.author_ms_department, ct.author_ms_gpa,
		ct.author_toefl, ct.author_ielts, ct.author_gre, ct.author_gmat, ct.author_testdaf, ct.author_goethe, ct.course_type
		from article as a		
		left join content as ct on (ct.id = a.id)
		where a.course_id= $1
		order by ct.published_date desc`

	articleRows, _ := m.DB.QueryContext(ctx, articleQuery, course.ID)
	// articles := make(map[int]Article)
	var articles []Article
	for articleRows.Next() {
		var ca CourseArticle
		err := articleRows.Scan(
			&ca.Article.ID,
			&ca.CourseID,
			&ca.Article.Result,
			&ca.Article.IsDecision,
			&ca.Article.Title,
			&ca.Article.Author,
			&ca.Article.Link,
			&ca.Article.PublishedAt,
			&ca.Article.Source,
			&ca.Article.AuthorBsSchool,
			&ca.Article.AuthorBsSchoolShort,
			&ca.Article.AuthorBsDepartment,
			&ca.Article.AuthorBsGpa,
			&ca.Article.AuthorMsSchool,
			&ca.Article.AuthorMsSchoolShort,
			&ca.Article.AuthorMsDepartment,
			&ca.Article.AuthorMsGpa,
			&ca.Article.AuthorToefl,
			&ca.Article.AuthorIelts,
			&ca.Article.AuthorGre,
			&ca.Article.AuthorGmat,
			&ca.Article.AuthorTestdaf,
			&ca.Article.AuthorGoethe,
			&ca.Article.CourseType,
		)
		if err != nil {
			return err
		}
		articles = append(articles, ca.Article)
		// articles[ca.ArticleID] = ca.Article
	}
	articleRows.Close()

	course.CourseArticle = articles
	return nil
}

func (m *DBModel) InsertCourse(course Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into course (id, university_id, course_type, name_en, name_en_short, 
		name_ch, name_ch_short, tuition_fees, beginning, subject, daadlink, is_elearning, application_deadline, is_complete_online_possible,
		programme_duration, is_from_daad, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`

	uid, _ := strconv.Atoi(course.UniversityId)
	ct, _ := strconv.Atoi(course.CourseType)

	_, err := m.DB.ExecContext(ctx, stmt,
		course.ID,
		uid,
		ct,
		course.NameEn,
		course.NameEnShort,
		course.NameCh,
		course.NameChShort,
		course.TuitionFees,
		course.Beginning,
		course.Subject,
		course.Daadlink,
		course.IsElearning,
		course.ApplicationDeadline,
		course.IsCompleteOnlinePossible,
		course.ProgrammeDuration,
		course.IsFromDaad,
		course.CreatedAt,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
