-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

create table users
(
    id           uuid primary key             default gen_random_uuid(),
    name         varchar(255)        not null,
    display_name varchar(255)        not null,
    email        varchar(255) unique not null,
    password     varchar(255)        not null,
    is_verified               not null default false,
    avatar_url   varchar(255),
    created_at   timestamp           not null default now(),
    edited_at    timestamp check (edited_at >= created_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
