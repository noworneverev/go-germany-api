package models

import (
	"context"
	"time"
)

func (m *DBModel) InsertUniversity(university University) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into university (id, name_en, name_ch, city, is_from_daad, is_tu9, is_u15, qs_ranking, created_at, updated_at, link) values 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := m.DB.ExecContext(ctx, stmt,
		university.ID,
		university.NameEn,
		university.NameCh,
		university.City,
		university.IsFromDaad,
		university.IsTu9,
		university.IsU15,
		university.QsRanking,
		university.CreatedAt,
		nil,
		university.Link,
	)
	if err != nil {
		return err
	}
	return nil
}
