CREATE TABLE IF NOT EXISTS users
(
    id serial primary key,
    name varchar(64),
    email varchar(255) not null unique,
    password varchar(255) not null,
    deleted_at timestamp default null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP
);
COMMENT ON COLUMN users.name IS '名前';
COMMENT ON COLUMN users.email IS 'メールアドレス';
COMMENT ON COLUMN users.password IS 'パスワード';

INSERT INTO users(
    name,
    email,
    password) values ('おかも', 'okmkm321@gmail.com', '9b7713f84cccd55251234fd5ae7905ad0748e009')