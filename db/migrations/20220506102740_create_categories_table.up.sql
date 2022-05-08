CREATE TABLE IF NOT EXISTS categories
(
    id serial primary key,
    name varchar(64) not null,
    slug varchar(64) not null unique,
    is_public boolean not null default false,
    position int unique,
    parent_id integer default null,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP,

    foreign key (parent_id) references categories(id)
    on delete SET NULL
);

COMMENT ON COLUMN categories.name IS 'カテゴリー名';
COMMENT ON COLUMN categories.slug IS 'スラッグ名';
COMMENT ON COLUMN categories.is_public IS '状態';
COMMENT ON COLUMN categories.position IS '順序';
COMMENT ON COLUMN categories.parent_id IS '親カテゴリ';