-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.events
(
    id serial NOT NULL,
    event_type text,
    event_date timestamp without time zone,
    accomplish_min_score integer,
    accomplish_max_score integer,
    expected_score integer,
    hook jsonb,
    expiration_date timestamp without time zone,
    remind_type text NOT NULL,
    remind_max_score integer,
    remind_description text NOT NULL,
    object_description text NOT NULL,
    remind_hook jsonb,
    CONSTRAINT pk_events PRIMARY KEY (id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.remind
(
    id serial NOT NULL,
    remind_description text,
    remind_type text NOT NULL,
    expire_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    status_description text,
    visibility text,
    event_id integer NOT NULL,
    max_score integer,
    object_description text,
    hook jsonb,
    CONSTRAINT reminder_pkey PRIMARY KEY (id),
    CONSTRAINT remind_events_fkey FOREIGN KEY (event_id)
        REFERENCES public.events (id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.accomplishers
(
    id serial not null,
    remind_id integer not null,
    event_id integer not null,
    accomplish_at timestamp without time zone,
    score integer not null,
    constraint accomplischers_pkey primary key (id),
    constraint accomplischers_reminders_fkey 
        foreign key (remind_id) 
        references remind(id),
    CONSTRAINT accomplischers_events_fkey 
        FOREIGN KEY (event_id)
        REFERENCES public.events (id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists public.accomplishers;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE if exists public.reminder;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE if exists public.events;
-- +goose StatementEnd