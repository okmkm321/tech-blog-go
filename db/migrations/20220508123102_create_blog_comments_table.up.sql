CREATE TABLE IF NOT EXISTS blog_comments
(
    id serial primary key,
    blog_id int not null,
    comment varchar not null,
    is_public boolean default true,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP,

    foreign key (blog_id) references blogs(id)
    on delete CASCADE
);

COMMENT ON COLUMN blog_comments.comment IS 'コメント';
