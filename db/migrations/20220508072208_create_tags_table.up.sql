CREATE TABLE IF NOT EXISTS tags
(
    id serial primary key,
    name varchar(64) not null unique,
    state integer not null default 0,
    position int unique,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP
);

COMMENT ON COLUMN tags.name IS 'タグ名';
COMMENT ON COLUMN tags.state IS '状態';
COMMENT ON COLUMN tags.position IS '順序';