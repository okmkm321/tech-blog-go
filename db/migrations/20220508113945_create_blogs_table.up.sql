CREATE TABLE IF NOT EXISTS blogs
(
    id serial primary key,
    user_id int not null,
    title varchar not null unique,
    description varchar,
    eye_catch varchar,
    body varchar,
    state int not null default 0,
    publish_at date,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP,

    foreign key (user_id) references users(id)
    on delete SET NULL
);

COMMENT ON COLUMN blogs.title IS 'タイトル';
COMMENT ON COLUMN blogs.description IS '見出し説明文';
COMMENT ON COLUMN blogs.eye_catch IS 'アイキャッチ';
COMMENT ON COLUMN blogs.body IS '本文';
COMMENT ON COLUMN blogs.state IS '状態';
COMMENT ON COLUMN blogs.publish_at IS '公開日';