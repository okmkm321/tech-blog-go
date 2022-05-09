package models

import (
	"context"
	"time"
)

type Blog struct {
	ID          int            `json:"id"`
	UserID      int            `json:"user_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	EyeCatch    string         `json:"eye_catch"`
	Body        string         `json:"body"`
	State       int            `json:"state"`
	PublishAt   string         `json:"publish_at"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	Category    []BlogCategory `json:"categories"`
	Tag         []BlogTag      `json:"tags"`
	Content     []BlogContent  `json:"contents"`
}

type BlogCategory struct {
	ID         int    `json:"-"`
	BlogID     int    `json:"-"`
	CategoryID int    `json:"-"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
}

type BlogTag struct {
	ID     int    `json:"-"`
	BlogID int    `json:"-"`
	TagID  int    `json:"-"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

type BlogContent struct {
	ID       int    `json:"-"`
	BlogID   int    `json:"-"`
	Name     string `json:"name"`
	Anchor   string `json:"anchor"`
	Position int    `json:"position"`
}

// all
func (m *DBModel) BlogGetAll() ([]*Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT
				id, title, description, eye_catch, body, state, publish_at
			FROM 
				blogs 
			WHERE
				blogs.deleted_at is null 
			ORDER BY publish_at ASC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bs []*Blog

	for rows.Next() {
		var b Blog
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.EyeCatch,
			&b.Body,
			&b.State,
			&b.PublishAt,
		)
		if err != nil {
			return nil, err
		}

		// category
		cQuery := `SELECT
					bc.id, bc.blog_id, bc.category_id, c.name, c.slug
				FROM
					blog_categories bc
					left join categories c on (c.id = bc.category_id)
				WHERE
					bc.blog_id = $1
				`

		cRows, _ := m.DB.QueryContext(ctx, cQuery, b.ID)

		if cRows != nil {
			var cs []BlogCategory
			for cRows.Next() {
				var bc BlogCategory
				err := cRows.Scan(
					&bc.ID,
					&bc.BlogID,
					&bc.CategoryID,
					&bc.Name,
					&bc.Slug,
				)
				if err != nil {
					return nil, err
				}
				cs = append(cs, bc)
			}
			cRows.Close()
			b.Category = cs
		}

		// tag
		tQuery := `SELECT
					bt.id, bt.blog_id, bt.tag_id, t.name, t.slug
				FROM
					blog_tags bt
					left join tags t on (t.id = bt.tag_id)
				WHERE
					bt.blog_id = $1
			`

		tRows, _ := m.DB.QueryContext(ctx, tQuery, b.ID)

		if tRows != nil {
			var ts []BlogTag
			for tRows.Next() {
				var bt BlogTag
				err := tRows.Scan(
					&bt.ID,
					&bt.BlogID,
					&bt.TagID,
					&bt.Name,
					&bt.Slug,
				)
				if err != nil {
					return nil, err
				}
				ts = append(ts, bt)
			}
			tRows.Close()
			b.Tag = ts
		}

		// content
		contentQuery := `SELECT
						id, blog_id, name, anchor, position
					FROM
						blog_contents
					WHERE
						blog_id = $1
					ORDER BY position ASC
					`

		contentRows, _ := m.DB.QueryContext(ctx, contentQuery, b.ID)
		if contentRows != nil {
			var contents []BlogContent
			for contentRows.Next() {
				var content BlogContent
				err := contentRows.Scan(
					&content.ID,
					&content.BlogID,
					&content.Name,
					&content.Anchor,
					&content.Position,
				)
				if err != nil {
					return nil, err
				}
				contents = append(contents, content)
			}
			contentRows.Close()
			b.Content = contents
		}

		bs = append(bs, &b)
	}

	return bs, nil
}

// GetOne
func (m *DBModel) GetBlog(id int) (*Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT
				id, title, description, eye_catch, body, state, publish_at
			FROM
				blogs
			WHERE
				id = $1
			ORDER BY publish_at ASC
			`
	row := m.DB.QueryRowContext(ctx, query, id)
	var b Blog
	err := row.Scan(
		&b.ID,
		&b.Title,
		&b.Description,
		&b.EyeCatch,
		&b.Body,
		&b.State,
		&b.PublishAt,
	)
	if err != nil {
		return nil, err
	}
	cQuery := `SELECT
				bc.id, bc.blog_id, bc.category_id, c.name, c.slug
			FROM
				blog_categories bc
				left join categories c on (c.id = bc.category_id)
			WHERE
				bc.blog_id = $1
	`
	cRows, _ := m.DB.QueryContext(ctx, cQuery, b.ID)

	if cRows != nil {
		var cs []BlogCategory
		for cRows.Next() {
			var bc BlogCategory
			err := cRows.Scan(
				&bc.ID,
				&bc.BlogID,
				&bc.CategoryID,
				&bc.Name,
				&bc.Slug,
			)
			if err != nil {
				return nil, err
			}
			cs = append(cs, bc)
		}
		cRows.Close()
		b.Category = cs
	}

	// tag
	tQuery := `SELECT
					bt.id, bt.blog_id, bt.tag_id, t.name, t.slug
				FROM
					blog_tags bt
					left join tags t on (t.id = bt.tag_id)
				WHERE
					bt.blog_id = $1
			`

	tRows, _ := m.DB.QueryContext(ctx, tQuery, b.ID)

	if tRows != nil {
		var ts []BlogTag
		for tRows.Next() {
			var bt BlogTag
			err := tRows.Scan(
				&bt.ID,
				&bt.BlogID,
				&bt.TagID,
				&bt.Name,
				&bt.Slug,
			)
			if err != nil {
				return nil, err
			}
			ts = append(ts, bt)
		}
		tRows.Close()
		b.Tag = ts
	}

	// content
	contentQuery := `SELECT
						id, blog_id, name, anchor, position
					FROM
						blog_contents
					WHERE
						blog_id = $1
					ORDER BY position ASC
					`

	contentRows, _ := m.DB.QueryContext(ctx, contentQuery, b.ID)
	if contentRows != nil {
		var contents []BlogContent
		for contentRows.Next() {
			var content BlogContent
			err := contentRows.Scan(
				&content.ID,
				&content.BlogID,
				&content.Name,
				&content.Anchor,
				&content.Position,
			)
			if err != nil {
				return nil, err
			}
			contents = append(contents, content)
		}
		contentRows.Close()
		b.Content = contents
	}

	return &b, nil
}

func (m *DBModel) BlogCreate(blog Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, _ := m.DB.Begin()
	defer func() {
		// panicが起きたらロールバック
		if recover() != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO blogs (
				user_id,
				title,
				description,
				eye_catch,
				body,
				state,
				publish_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7
			) RETURNING id`
	lid := 0
	err := tx.QueryRowContext(ctx, query,
		blog.UserID,
		blog.Title,
		blog.Description,
		blog.EyeCatch,
		blog.Body,
		blog.State,
		blog.PublishAt,
	).Scan(&lid)
	if err != nil {
		tx.Rollback() // ロールバック
		return err
	}

	// カテゴリー登録
	for i := 0; i < len(blog.Category); i++ {
		cQuery := `INSERT INTO blog_categories (
					blog_id,
					category_id
				) VALUES (
					$1, $2
				)`

		_, err = tx.ExecContext(ctx, cQuery, lid, blog.Category[i].CategoryID)
		if err != nil {
			tx.Rollback() // ロールバック
			return err
		}
	}

	// タグ登録
	for i := 0; i < len(blog.Tag); i++ {
		tQuery := `INSERT INTO blog_tags (
					blog_id,
					tag_id
				) VALUES (
					$1, $2
				)`

		_, err = tx.ExecContext(ctx, tQuery, lid, blog.Tag[i].TagID)
		if err != nil {
			tx.Rollback() // ロールバック
			return err
		}
	}

	for i := 0; i < len(blog.Content); i++ {
		cQuery := `INSERT INTO blog_contents (
					blog_id,
					name,
					anchor,
					position
				) VALUES (
					$1, $2, $3, $4
				)`

		_, err = tx.ExecContext(ctx, cQuery, lid, blog.Content[i].Name, blog.Content[i].Anchor, blog.Content[i].Position)
		if err != nil {
			tx.Rollback() // ロールバック
			return err
		}
	}
	tx.Commit()
	return nil

}
