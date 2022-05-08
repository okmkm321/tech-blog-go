CREATE TABLE IF NOT EXISTS blog_images
(
    id serial primary key,
    blog_id int not null,
    path varchar not null,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP,

    foreign key (blog_id) references blogs(id)
    on delete CASCADE
);

COMMENT ON COLUMN blog_images.path IS '画像パス';
