create table try_outs
(
    id           text not null
        primary key,
    created_at   bigint,
    updated_at   bigint,
    deleted_at   timestamp with time zone,
    tryout_name  text,
    tryout_type  text,
    tryout_quota integer,
    version      bigint,
    tryout_price integer,
    due_date     bigint
);

alter table try_outs
    owner to postgres;

create table participants
(
    id                varchar(255) not null
        primary key,
    created_at        bigint,
    updated_at        bigint,
    deleted_at        timestamp,
    participants_name varchar(255),
    version           integer,
    class             varchar(255),
    city              varchar(255),
    province          varchar(255),
    whatsapp          varchar(255),
    email             varchar(255),
    token             varchar(20)
        constraint participants_pk
            unique
);

alter table participants
    owner to postgres;

create table submission
(
    id                    text not null
        primary key,
    created_at            bigint,
    updated_at            bigint,
    deleted_at            timestamp with time zone,
    participants_id       text
        constraint submission_participants_id_fk
            references participants,
    tryout_id             text
        references try_outs,
    submission_start_time timestamp with time zone,
    version               bigint,
    submission_end_time   timestamp with time zone,
    token                 varchar(255)
        constraint submission_participants_token_fk
            references participants (token)
);

alter table submission
    owner to postgres;

create table answer_submitted
(
    id                      text not null
        primary key,
    created_at              bigint,
    updated_at              bigint,
    deleted_at              timestamp with time zone,
    user_test_submission_id text
        references submission,
    question_id             text,
    choice_id               text,
    version                 bigint
);

alter table answer_submitted
    owner to postgres;

