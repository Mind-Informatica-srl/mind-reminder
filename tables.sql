-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.remind_to_calculate
(
    id serial NOT NULL,
    action text not null,
    object_id text not null,
    object_type text not null,
    object_raw jsonb,
    created_at timestamp without time zone DEFAULT now(),
    elaborated_at timestamp without time zone,
    error text,
    CONSTRAINT to_remind_pkey PRIMARY KEY
(id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.remind
(
    id serial NOT NULL,
    description text,
    remind_type text not null,
    object_raw jsonb,
    object_type text not null,
    expire_at timestamp without time zone,
    accomplished_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    percentage numeric,
    status_description text,
    visibility text,
    object_id text,
    CONSTRAINT reminder_pkey PRIMARY KEY
(id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.accomplishers
(
    id serial not null,
    remind_id integer not null,
    object_id text not null,
    accomplish_at timestamp without time zone,
    percetange FLOAT,
    constraint accomplischers_pkey primary key (id),
    constraint accomplischers_reminders_fkey foreign key (remind_id) references remind(id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists public.accomplishers;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE if exists public.remind_to_calculate;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE if exists public.reminder;
-- +goose StatementEnd