package models

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) Get(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, email from users where id = $1`
	fmt.Println(query)

	row := m.DB.QueryRowContext(ctx, query, id)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *DBModel) All() ([]*User, error) {
	return nil, nil
}
