create table choices
(
    id         text not null
        primary key,
    created_at bigint,
    updated_at bigint,
    content    text,
    is_correct boolean,
    weight     integer,
    version    integer
);

alter table choices
    owner to postgres;

create table submission
(
    id         text not null
        primary key,
    created_at bigint,
    updated_at bigint,
    token      varchar(20),
    score      integer,
    version    integer
);

alter table submission
    owner to postgres;

create table answer_submitted
(
    id                      text not null
        primary key,
    created_at              bigint,
    updated_at              bigint,
    user_test_submission_id text
        references submission,
    choice_id               text
        references choices,
    version                 integer
);

alter table answer_submitted
    owner to postgres;

