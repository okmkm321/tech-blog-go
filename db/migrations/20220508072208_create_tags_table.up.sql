CREATE TABLE IF NOT EXISTS tags
(
    id serial primary key,
    name varchar(64) not null unique,
    is_public boolean not null default false,
    position int unique,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP
);

COMMENT ON COLUMN tags.name IS 'タグ名';
COMMENT ON COLUMN tags.is_public IS '状態';
COMMENT ON COLUMN tags.position IS '順序';