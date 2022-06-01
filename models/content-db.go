package models

import (
	"context"
	"time"
)

func (m *DBModel) InsertContent(content Content) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into content (id, link, title, author, published_date, source, 
		author_bs_school, author_bs_school_short, author_bs_department, author_bs_gpa,
		author_ms_school, author_ms_school_short, author_ms_department, author_ms_gpa,
		author_toefl, author_ielts, author_gre, author_gmat, author_testdaf, author_goethe, course_type)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`

	_, err := m.DB.ExecContext(ctx, stmt,
		content.ID,
		content.Link,
		content.Title,
		content.Author,
		content.PublishedAt,
		content.Source,
		content.AuthorBsSchool,
		content.AuthorBsSchoolShort,
		content.AuthorBsDepartment,
		content.AuthorBsGpa,
		content.AuthorMsSchool,
		content.AuthorMsSchoolShort,
		content.AuthorMsDepartment,
		content.AuthorMsGpa,
		content.AuthorToefl,
		content.AuthorIelts,
		content.AuthorGre,
		content.AuthorGmat,
		content.AuthorTestdaf,
		content.AuthorGoethe,
		content.CourseType,
	)

	if err != nil {
		return err
	}
	return nil
}
