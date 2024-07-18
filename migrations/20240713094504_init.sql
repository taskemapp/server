-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

create table users
(
    id           uuid primary key             default gen_random_uuid(),
    name         varchar(255) unique not null check (length(name) > 3),
    display_name varchar(255)        not null check (length(display_name) > 3),
    email        varchar(255) unique not null,
    password     varchar(255)        not null,
    is_verified  bool                    not null default false,
    avatar_url   varchar(255),
    created_at   timestamp           not null default now(),
    edited_at    timestamp check (edited_at >= created_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
