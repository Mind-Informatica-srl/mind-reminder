-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.to_remind
(
    id serial NOT NULL,
    action text not null,
    object_id text not null,
    object_type text not null,
    object_raw json,
    created_at timestamp without time zone DEFAULT now(),
    elaborated_at timestamp without time zone,
    error text,
    CONSTRAINT to_remind_pkey PRIMARY KEY
(id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.reminder
(
    id serial NOT NULL,
    description text,
    reminder_type text not null,
    object_raw json,
    object_type text not null,
    expire_at timestamp without time zone,
    accomplished_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    percentage numeric,
    status_description text,
    visibility text,
    CONSTRAINT reminder_pkey PRIMARY KEY
(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists public.to_remind;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE if exists public.reminder;
-- +goose StatementEnd