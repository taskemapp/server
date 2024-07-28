-- +goose Up
-- +goose StatementBegin
create table teams
(
    id               uuid primary key             default gen_random_uuid(),
    name             varchar(255) unique not null check (length(name) > 3),
    description      varchar(255)        not null check (length(description) > 3),
    header_image_Url varchar(255) unique not null,
    creator          uuid                not null references users (id) on delete cascade,
    created_at       timestamp           not null default now(),
    edited_at        timestamp check (edited_at >= created_at)
);
-- +goose StatementEnd

-- +goose StatementBegin
create table team_members
(
    id        uuid primary key   default gen_random_uuid(),
    user_id   uuid      not null references users (id) on delete cascade,
    team_id   uuid      not null references teams (id) on delete cascade,
    joined_at timestamp not null default now(),
    leaved_at timestamp check (leaved_at >= joined_at),
    is_leaved bool      not null default false,
    constraint unique_user_team_combination unique (user_id, team_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
