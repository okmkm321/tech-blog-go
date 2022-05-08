CREATE TABLE IF NOT EXISTS blog_contents
(
    id serial primary key,
    blog_id int not null,
    name varchar not null,
    anchor varchar not null,
    position int,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP,

    foreign key (blog_id) references blogs(id)
    on delete CASCADE
);

COMMENT ON COLUMN blog_contents.name IS '目次名';
COMMENT ON COLUMN blog_contents.anchor IS 'アンカー名';
COMMENT ON COLUMN blog_contents.position IS '順序';
