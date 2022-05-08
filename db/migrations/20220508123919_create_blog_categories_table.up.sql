CREATE TABLE IF NOT EXISTS blog_categories
(
    id serial primary key,
    blog_id int not null,
    category_id int not null,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP,

    foreign key (blog_id) references blogs(id)
    on delete CASCADE,

    foreign key (category_id) references categories(id)
    on delete CASCADE
);