-- +goose Up
create table users (
    id serial primary key,
    username text unique not null,
    password text not null,
    role text not null,
    created_at timestamp not null,
    updated_at timestamp
);

-- +goose Down
drop table users;
