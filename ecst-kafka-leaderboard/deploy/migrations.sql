create table participants
(
    id                text not null
        primary key,
    created_at        bigint,
    updated_at        bigint,
    participants_name text,
    version           bigint,
    token             varchar(20)
        constraint participants_pk
            unique
);

alter table participants
    owner to postgres;

create table grades
(
    id         text not null
        primary key,
    created_at bigint,
    updated_at bigint,
    tryout_id  text,
    score      integer,
    version    bigint,
    token      varchar(20)
        constraint grades_participants_token_fk
            references participants (token)
);

alter table grades
    owner to postgres;

