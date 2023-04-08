create table if not exists permissions
(
    id           int auto_increment
    primary key,
    name         varchar(50)                        not null,
    tenant_id    int                                not null,
    created_time datetime default CURRENT_TIMESTAMP not null,
    updated_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
    );

create table if not exists role_permissions
(
    id            int auto_increment
    primary key,
    role_id       int                                not null,
    permission_id int                                not null,
    created_time  datetime default CURRENT_TIMESTAMP not null,
    updated_time  datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint role_id
    unique (role_id, permission_id)
    );

create table if not exists roles
(
    id           int auto_increment
    primary key,
    name         varchar(50)                        not null,
    tenant_id    int                                not null,
    created_time datetime default CURRENT_TIMESTAMP not null,
    updated_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
    );

create table if not exists tenant
(
    id           bigint(32) auto_increment
    primary key,
    name         varchar(32)                        not null,
    created_time datetime default CURRENT_TIMESTAMP not null,
    updated_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    description  text                               null
    );

create table if not exists third_auth
(
    id              int auto_increment
    primary key,
    user_id         int                                not null,
    third_auth_type int                                not null,
    third_auth_id   varchar(100)                       not null,
    token           varchar(255)                       not null,
    expires_time    datetime                           null,
    created_time    datetime default CURRENT_TIMESTAMP not null,
    updated_time    datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint third_auth_type
    unique (third_auth_type, third_auth_id)
    );

create table if not exists users
(
    id           int auto_increment
    primary key,
    email        varchar(50)                        not null,
    nick_name    varchar(50)                        not null,
    password     varchar(1024)                        not null,
    tenant_id    int                                not null,
    created_time datetime default CURRENT_TIMESTAMP not null,
    updated_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
    );

create table if not exists verification_code
(
    id                int auto_increment
    primary key,
    send_subject_type int                                not null,
    send_subject_name varchar(255)                       not null,
    code              varchar(50)                        not null,
    expire_time       datetime                           not null,
    created_time      datetime default CURRENT_TIMESTAMP not null,
    updated_time      datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
    );

