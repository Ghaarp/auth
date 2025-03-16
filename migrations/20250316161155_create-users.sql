-- +goose Up
create table users (
    id serial primary key,
    username text not null,
    email text not null,
    pass_hash text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- +goose Down
drop table users;
