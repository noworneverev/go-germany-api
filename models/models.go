package models

import (
	"database/sql"
	"time"
)

// Models is the wrapper for database
type Models struct {
	DB DBModel
}

// NewModles returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Course is the type for courses
type Course struct {
	ID                       int       `json:"id"`
	UniversityId             string    `json:"university_id"`
	CourseType               string    `json:"course_type"`
	NameEn                   string    `json:"name_en"`
	NameEnShort              string    `json:"name_en_short"`
	NameCh                   string    `json:"-"`
	NameChShort              string    `json:"-"`
	TuitionFees              string    `json:"tuition_fees"`
	Beginning                string    `json:"beginning"`
	Subject                  string    `json:"subject"`
	Daadlink                 string    `json:"daadlink"`
	IsElearning              bool      `json:"is_elearning"`
	ApplicationDeadline      string    `json:"application_deadline"`
	IsCompleteOnlinePossible bool      `json:"is_complete_online_possible"`
	ProgrammeDuration        string    `json:"programme_duration"`
	IsFromDaad               bool      `json:"is_from_daad"`
	CreatedAt                time.Time `json:"-"`
	UpdatedAt                time.Time `json:"-"`
	UniversityNameEn         string    `json:"university_name_en"`
	UniversityNameCh         string    `json:"university_name_ch"`
	City                     string    `json:"city"`
	IsTu9                    bool      `json:"is_tu9"`
	IsU15                    bool      `json:"is_u15"`
	QsRanking                int       `json:"qs_ranking"`
	UniversityLink           string    `json:"university_link"`
	// CourseLanguage           map[int]string `json:"languages"`
	// CourseArticle            map[int]Article `json:"articles"`
	CourseLanguage []string  `json:"languages"`
	CourseArticle  []Article `json:"articles"`
	Languages      string    `json:"-"`
	ArticleCount   int       `json:"-"`
}

// Language is the type for languages
type Language struct {
	ID           int    `json:"id"`
	LanguageName string `json:"language_name"`
}

// CourseLanguage is the type for course language
type CourseLanguage struct {
	ID         int      `json:"-"`
	CourseID   int      `json:"-"`
	LanguageID int      `json:"-"`
	Language   Language `json:"language"`
}

// Article is the type for articles
type Article struct {
	ID                  int             `json:"id"`
	Title               string          `json:"title"`
	Author              string          `json:"author"`
	Link                string          `json:"link"`
	PublishedAt         time.Time       `json:"published_at"`
	Source              string          `json:"source"`
	AuthorBsSchool      string          `json:"author_bs_school"`
	AuthorBsSchoolShort string          `json:"author_bs_school_short"`
	AuthorBsDepartment  string          `json:"author_bs_department"`
	AuthorBsGpa         string          `json:"author_bs_gpa"`
	AuthorMsSchool      string          `json:"author_ms_school"`
	AuthorMsSchoolShort string          `json:"author_ms_school_short"`
	AuthorMsDepartment  string          `json:"author_ms_department"`
	AuthorMsGpa         string          `json:"author_ms_gpa"`
	AuthorToefl         string          `json:"author_toefl"`
	AuthorIelts         string          `json:"author_ielts"`
	AuthorGre           string          `json:"author_gre"`
	AuthorGmat          string          `json:"author_gmat"`
	AuthorTestdaf       string          `json:"author_testdaf"`
	AuthorGoethe        string          `json:"author_goethe"`
	CourseType          string          `json:"course_type"`
	Result              string          `json:"result"`
	IsDecision          bool            `json:"is_decision"`
	ArticleCourse       []ArticleCourse `json:"courses"`
	Content             string          `json:"content"`
}

// CourseArticle is the type for course article
type CourseArticle struct {
	ArticleID  int     `json:"-"`
	CourseID   int     `json:"-"`
	Result     string  `json:"result"`
	IsDecision bool    `json:"is_decision"`
	Article    Article `json:"article"`
}

type CourseParams struct {
	PageNumber           int    `json:"page_number"`
	PageSize             int    `json:"page_size"`
	SearchTerm           string `json:"search_term"`
	CourseTypes          string `json:"course_types"`
	Languages            string `json:"languages"`
	Subjects             string `json:"subjects"`
	Institutions         string `json:"institutions"`
	IsTu9                bool   `json:"is_tu9"`
	IsU15                bool   `json:"is_u15"`
	HasArticles          bool   `json:"has_articles"`
	OrderBy              string `json:"orderby"`
	HideLanguageNArticle bool   `json:"hide_language_article"`
}

type Filters struct {
	CourseTypes  []string `json:"course_types"`
	Languages    []string `json:"languages"`
	Subjects     []string `json:"subjects"`
	Institutions []string `json:"institutions"`
}

type ArticleCourse struct {
	ArticleID  int    `json:"-"`
	CourseID   int    `json:"-"`
	Result     string `json:"result"`
	IsDecision bool   `json:"is_decision"`
	Course     Course `json:"course"`
}

type ArticleParams struct {
	PageNumber      int    `json:"page_number"`
	PageSize        int    `json:"page_size"`
	SearchTerm      string `json:"search_term"`
	Sources         string `json:"sources"`
	BsSchools       string `json:"bs_schools"`
	BsDepartments   string `json:"bs_departments"`
	MsSchools       string `json:"ms_schools"`
	MsDepartments   string `json:"ms_departments"`
	CourseType      string `json:"course_type"`
	HideApplication bool   `json:"hide_application"`
}

type ArticleFilters struct {
	Sources       []string `json:"sources"`
	BsSchools     []string `json:"bs_schools"`
	BsDepartments []string `json:"bs_departments"`
	MsSchools     []string `json:"ms_schools"`
	MsDepartments []string `json:"ms_departments"`
	CourseTypes   []string `json:"course_types"`
}

// User is the type for users
type User struct {
	ID       int
	Email    string
	Password string
}

// University is the type for universities
type University struct {
	ID         int       `json:"id"`
	NameEn     string    `json:"name_en"`
	NameCh     string    `json:"-"`
	City       string    `json:"city"`
	IsFromDaad bool      `json:"is_from_daad"`
	IsTu9      bool      `json:"is_tu9"`
	IsU15      bool      `json:"is_u15"`
	QsRanking  int       `json:"qs_ranking"`
	CreatedAt  time.Time `json:"-"`
	Link       string    `json:"link"`
	UpdatedAt  time.Time `json:"-"`
}

// Content is the type for content
type Content struct {
	ID                  int       `json:"id"`
	Link                string    `json:"link"`
	Title               string    `json:"title"`
	Author              string    `json:"author"`
	PublishedAt         time.Time `json:"published_at"`
	Source              string    `json:"source"`
	AuthorBsSchool      string    `json:"author_bs_school"`
	AuthorBsSchoolShort string    `json:"author_bs_school_short"`
	AuthorBsDepartment  string    `json:"author_bs_department"`
	AuthorBsGpa         string    `json:"author_bs_gpa"`
	AuthorMsSchool      string    `json:"author_ms_school"`
	AuthorMsSchoolShort string    `json:"author_ms_school_short"`
	AuthorMsDepartment  string    `json:"author_ms_department"`
	AuthorMsGpa         string    `json:"author_ms_gpa"`
	AuthorToefl         string    `json:"author_toefl"`
	AuthorIelts         string    `json:"author_ielts"`
	AuthorGre           string    `json:"author_gre"`
	AuthorGmat          string    `json:"author_gmat"`
	AuthorTestdaf       string    `json:"author_testdaf"`
	AuthorGoethe        string    `json:"author_goethe"`
	CourseType          string    `json:"course_type"`
	Content             string    `json:"Content"`
}
