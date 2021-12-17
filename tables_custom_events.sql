-- +goose Up
-- +goose StatementBegin
create table custom_object_prototypes(
    id serial not null,
    name varchar(100),
    description text,
    prototype_object_data jsonb,
    description_template text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    constraint pk_custom_object_prototypes primary key (id)
);
-- +goose StatementEnd
-- +goose StatementBegin
create table custom_event_prototypes(
    id serial not null,
    name varchar(100),
    description text,
    prototype_event_data jsonb,
    event_type_key text,
    event_date_key text,
    accomplishable_after_remind boolean default false,
    required_on_object_creation boolean default false,
    accomplish_min_score_key text,
    accomplish_max_score_key text,
    expected_score_key text,
    hook_keys text [],
    remind_expiration_date_key text,
    remind_type_key text,
    remind_max_score_key text,
    remind_description_template text,
    remind_object_description_template text,
    remind_hook_keys text [],
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    constraint pk_custom_event_prototypes primary key (id)
);
-- +goose StatementEnd
-- +goose StatementBegin
create table custom_events(
    id serial not null,
    custom_event_prototype_id integer not null,
    data jsonb,
    custom_object_id integer,
    constraint pk_custom_events primary key (id),
    CONSTRAINT fk_custom_ev_custom_ev_prototypes FOREIGN KEY (custom_event_prototype_id) REFERENCES public.custom_event_prototypes (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
create table custom_objects(
    id serial not null,
    custom_object_prototype_id integer not null,
    data jsonb,
    constraint pk_custom_objects primary key (id),
    CONSTRAINT fk_custom_ob_custom_ob_prototypes FOREIGN KEY (custom_object_prototype_id) REFERENCES public.custom_object_prototypes (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
create table custom_sections(
    id serial not null,
    name varchar(100),
    description text,
    icon varchar(50),
    configuration text,
    section_order integer default 0,
    reference text,
    custom_object_prototype_id integer,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    constraint pk_custom_section primary key (id),
    CONSTRAINT fk_custom_sections_custom_object_prototypes FOREIGN KEY (custom_object_prototype_id) REFERENCES public.custom_object_prototypes (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
create table public.custom_sections_custom_event_prototypes(
    custom_sections_id integer not null,
    custom_event_prototype_id integer,
    constraint pk_custom_sec_custom_ev_prototypes primary key (custom_sections_id, custom_event_prototype_id),
    CONSTRAINT fk_custom_sec_custom_ev_prototypes_custom_sections FOREIGN KEY (custom_sections_id) REFERENCES public.custom_sections (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_custom_sec_custom_ev_prototypes_custom_ev_prototypes FOREIGN KEY (custom_event_prototype_id) REFERENCES public.custom_event_prototypes (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists custom_sections_custom_event_prototypes;
-- +goose StatementEnd
-- +goose StatementBegin
drop table if exists custom_sections;
-- +goose StatementEnd
-- +goose StatementBegin
drop table if exists custom_objects;
-- +goose StatementEnd
-- +goose StatementBegin
drop table if exists custom_events;
-- +goose StatementEnd
-- +goose StatementBegin
drop table if exists custom_object_prototypes;
-- +goose StatementEnd
-- +goose StatementBegin
drop table if exists custom_event_prototypes;
-- +goose StatementEnd