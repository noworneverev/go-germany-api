package models

import (
	"context"
	"time"
)

func (m *DBModel) GetUser(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, email, password from gogermany_user where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
