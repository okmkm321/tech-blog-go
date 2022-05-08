package models

import (
	"context"
	"time"
)

type Category struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	IsPublic  bool        `json:"is_public"`
	Position  int         `json:"position"`
	ParentId  interface{} `json:"parent_id"`
	UpdatedAt time.Time   `json:"-"`
}

// All
func (m *DBModel) CategoryGetAll() ([]*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, slug, is_public, position, parent_id from categories where categories.deleted_at is null order by position ASC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		var ctg Category
		err := rows.Scan(
			&ctg.ID,
			&ctg.Name,
			&ctg.Slug,
			&ctg.IsPublic,
			&ctg.Position,
			&ctg.ParentId,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &ctg)
	}
	return categories, nil
}

// getOne
func (m *DBModel) GetCategory(id int) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, slug, position, parent_id from categories where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var ctg Category
	err := row.Scan(
		&ctg.ID,
		&ctg.Name,
		&ctg.Slug,
		&ctg.IsPublic,
		&ctg.Position,
		&ctg.ParentId,
	)
	if err != nil {
		return nil, err
	}
	return &ctg, nil
}

// Create
func (m *DBModel) CategoryCreate(category Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var pi interface{}
	if category.ParentId == 0 {
		pi = nil
	} else {
		pi = category.ParentId
	}

	count, err := m.countCategoryAll(category)
	if err != nil {
		return err
	}
	query := `insert into categories (name, slug, is_public, position, parent_id) values ($1, $2, $3, $4, $5)`

	_, err = m.DB.ExecContext(ctx, query,
		category.Name,
		category.Slug,
		category.IsPublic,
		count,
		pi,
	)

	if err != nil {
		return err
	}

	return nil
}

// UPDATE
func (m *DBModel) CategoryUpdate(category Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update categories set name = $1, slug = $2, is_public = $3, parent_id = $4, updated_at = $5 where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		category.Name,
		category.Slug,
		category.IsPublic,
		category.ParentId,
		category.UpdatedAt,
		category.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// delete
func (m *DBModel) CategoryDelete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, _ := m.DB.Begin()
	defer func() {
		// panicが起きたらロールバック
		if recover() != nil {
			tx.Rollback()
		}
	}()

	query := `select position from categories where id = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var ctg Category
	err := row.Scan(&ctg.Position)
	if err != nil {
		tx.Rollback() // ロールバック
		return err
	}

	query = `update categories set position = null, deleted_at = NOW() where id = $1`
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback() // ロールバック
		return err
	}

	query = `update categories set position = position - 1 where position > $1`
	_, err = tx.ExecContext(ctx, query, ctg.Position)
	if err != nil {
		tx.Rollback() // ロールバック
		return err
	}
	tx.Commit()

	return nil
}

// Record Count
func (m *DBModel) countCategoryAll(category Category) (count int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(*) from categories where categories.deleted_at is null`
	row := m.DB.QueryRowContext(ctx, query)
	err = row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count + 1, nil
}
