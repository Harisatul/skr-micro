
create table try_outs
(
    id           varchar(255) not null
        primary key,
    created_at   bigint,
    updated_at   bigint,
    deleted_at   timestamp,
    tryout_name  varchar(255),
    tryout_type  varchar(255),
    tryout_quota integer,
    version      bigint,
    tryout_price integer,
    due_date     bigint
);

alter table try_outs
    owner to postgres;

create table questions
(
    id         text not null
        primary key,
    created_at bigint,
    updated_at bigint,
    deleted_at timestamp with time zone,
    tryout_id  text
        references try_outs,
    content    text,
    weight     integer,
    version    bigint
);

alter table questions
    owner to postgres;

create table choices
(
    id          text not null
        primary key,
    created_at  bigint,
    updated_at  bigint,
    deleted_at  timestamp with time zone,
    question_id text
        references questions,
    content     text,
    is_correct  boolean,
    weight      integer,
    version     bigint
);

alter table choices
    owner to postgres;

create table users
(
    username   varchar(50)  not null
        constraint users_pk
            primary key,
    password   varchar(255) not null,
    name       varchar(50),
    email      varchar(255),
    whatsapp   varchar(255),
    token      text,
    created_at bigint,
    updated_at bigint,
    version    integer
);

alter table users
    owner to postgres;
