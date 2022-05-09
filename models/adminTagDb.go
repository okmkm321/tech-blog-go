package models

import (
	"context"
	"time"
)

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	IsPublic  bool      `json:"is_public"`
	Position  int       `json:"position"`
	UpdatedAt time.Time `json:"-"`
}

// all
func (m *DBModel) TagGetAll() ([]*Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, slug, is_public, position from tags where tags.deleted_at is null order by position ASC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug,
			&tag.IsPublic,
			&tag.Position,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

// getOne
func (m *DBModel) GetTag(id int) (*Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, slug, is_public, position from tags where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var t Tag
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Slug,
		&t.IsPublic,
		&t.Position,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Create
func (m *DBModel) TagCreate(tag Tag) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	count, err := m.countTagAll(tag)
	if err != nil {
		return err
	}
	query := `insert into tags (name, slug, is_public, position) values ($1, $2, $3, $4)`

	_, err = m.DB.ExecContext(ctx, query,
		tag.Name,
		tag.Slug,
		tag.IsPublic,
		count,
	)

	if err != nil {
		return err
	}

	return nil
}

// update
func (m *DBModel) TagUpdate(tag Tag) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update tags set name = $1, slug = $2, is_public = $3, updated_at = $4 where id = $5`

	_, err := m.DB.ExecContext(ctx, query,
		tag.Name,
		tag.Slug,
		tag.IsPublic,
		tag.UpdatedAt,
		tag.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// 削除
func (m *DBModel) TagDelete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, _ := m.DB.Begin()
	defer func() {
		// panicが起きたらロールバック
		if recover() != nil {
			tx.Rollback()
		}
	}()

	query := `select position from tags where id = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var t Tag
	err := row.Scan(&t.Position)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `update tags set position = null, deleted_at = NOW() where id = $1`
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback() // ロールバック
		return err
	}

	query = `update tags set position = position - 1 where position > $1`
	_, err = tx.ExecContext(ctx, query, t.Position)
	if err != nil {
		tx.Rollback() // ロールバック
		return err
	}
	tx.Commit()

	return nil
}

// Record Count
func (m *DBModel) countTagAll(tag Tag) (count int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(*) from tags where tags.deleted_at is null`
	row := m.DB.QueryRowContext(ctx, query)
	err = row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count + 1, nil
}
