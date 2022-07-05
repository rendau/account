do
$$
    begin
        execute 'ALTER DATABASE ' || current_database() || ' SET timezone = ''+06''';
    end;
$$;

create table cfg
(
    v jsonb not null default '{}'
);

create table role
(
    id        text not null
        primary key,
    name      text not null default '',
    is_system bool not null default false
);

create table perm
(
    id        text not null
        primary key,
    app       text not null default '',
    dsc       text not null default '',
    is_system bool not null default false
);

create table role_perm
(
    role_id text not null
        constraint role_perm_fk_role_id references role (id) on update cascade on delete cascade,
    perm_id text not null
        constraint role_perm_fk_perm_id references perm (id) on update cascade on delete cascade
);

create table usr
(
    id         bigserial   not null
        primary key,
    created_at timestamptz not null default now(),
    phone      text        not null
        constraint usr_unique_phone unique,
    ava        text        not null default '',
    name       text        not null default '',
    token      text        not null default ''
);
create index usr_created_at_idx
    on usr (created_at);
create index usr_phone_idx
    on usr (phone);
create index usr_token_idx
    on usr (token);

create table usr_role
(
    usr_id  bigint not null
        constraint usr_role_fk_usr_id references usr (id) on update cascade on delete cascade,
    role_id text   not null
        constraint usr_role_fk_role_id references role (id) on update cascade on delete cascade
);
create index usr_role_usr_id_idx
    on usr_role (usr_id);

do
$$
    declare
    begin
        -- Admin role
        insert into role(id, name, is_system)
        values ('admin', 'Admin', true);

        -- perms
        insert into perm(id, dsc, is_system)
        values ('*', 'All permissions', true)
             , ('m_perm', 'Modify permissions', true)
             , ('m_role', 'Modify roles', true)
             , ('m_usr', 'Modify users', true);

        -- Admin role_perm
        insert into role_perm(role_id, perm_id)
        values ('admin', '*');

        -- Admin user
        insert into usr(phone, name)
        values ('70000000000', 'Admin');

        -- Admin usr_role
        insert into usr_role(usr_id, role_id)
        select u.id, 'admin'
        from usr u
        where u.phone = '70000000000';
    end ;
$$;
