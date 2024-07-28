-- +goose Up
-- +goose StatementBegin

create type task_status AS ENUM (
    'in_progress',
    'paused',
    'finished',
    'canceled',
    'awaiting_assignment',
    'backlog_assigned'
    );
-- +goose StatementEnd

-- +goose StatementBegin
create type language_code AS ENUM ('en', 'ru');
-- +goose StatementEnd

-- +goose StatementBegin
create table task_status_translations
(
    status        task_status   not null,
    language_code language_code not null,
    translation   varchar(255),
    primary key (status, language_code),
    constraint unique_status_lang_combination unique (status, language_code)
);
-- +goose StatementEnd

-- +goose StatementBegin
insert into task_status_translations (status, language_code, translation)
values ('in_progress', 'en', 'In Progress'),
       ('in_progress', 'ru', 'В процессе выполнения'),
       ('paused', 'en', 'Paused'),
       ('paused', 'ru', 'Приостановлена'),
       ('finished', 'en', 'Finished'),
       ('finished', 'ru', 'Завершена'),
       ('canceled', 'en', 'Canceled'),
       ('canceled', 'ru', 'Отменена'),
       ('awaiting_assignment', 'en', 'Awaiting Assignment'),
       ('awaiting_assignment', 'ru', 'Ожидание назначения исполнителя'),
       ('backlog_assigned', 'en', 'Backlog Assigned'),
       ('backlog_assigned', 'ru', 'Назначен исполнитель, находится в бэклоге');
-- +goose StatementEnd

-- +goose StatementBegin
create table tasks
(
    id          uuid primary key      default gen_random_uuid(),
    name        varchar(255) not null check ( length(name) >= 3 ),
    description text         not null,
    status      task_status  not null,
    team_id     uuid         not null references teams (id) on delete cascade,
    assigned_id uuid references users (id) on delete cascade,
    creator     uuid         not null references users (id) on delete cascade,
    created_at  timestamp    not null default now(),
    edited_at   timestamp check (edited_at >= created_at),
    ended_at    timestamp check (ended_at >= created_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
