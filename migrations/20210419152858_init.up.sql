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

create table app
(
    id             bigserial not null
        primary key,
    name           text      not null default '',
    perm_url       text      not null default '',
    is_account_app bool      not null default false
);

create table perm
(
    id        bigserial not null
        primary key,
    code      text      not null default '',
    is_all    bool      not null default false,
    app_id    bigint    not null
        constraint perm_fk_app_id references app (id) on update cascade on delete cascade,
    dsc       text      not null default '',
    is_system bool      not null default false
);

create table role
(
    id        bigserial not null
        primary key,
    code      text      not null default ''
        constraint role_unique_code unique,
    name      text      not null default '',
    is_system bool      not null default false
);

create table role_perm
(
    role_id bigint not null
        constraint role_perm_fk_role_id references role (id) on update cascade on delete cascade,
    perm_id bigint not null
        constraint role_perm_fk_perm_id references perm (id) on update cascade on delete cascade
);

create table usr
(
    id         bigserial   not null
        primary key,
    created_at timestamptz not null default now(),
    active     bool        not null default true,
    phone      text        not null
        constraint usr_unique_phone unique,
    ava        text        not null default '',
    name       text        not null default ''
);
create index usr_created_at_idx
    on usr (created_at);
create index usr_active_idx
    on usr (active);
create index usr_phone_idx
    on usr (phone);

create table usr_role
(
    usr_id  bigint not null
        constraint usr_role_fk_usr_id references usr (id) on update cascade on delete cascade,
    role_id bigint not null
        constraint usr_role_fk_role_id references role (id) on update cascade on delete cascade
);
create index usr_role_usr_id_idx
    on usr_role (usr_id);

do
$$
    declare
        account_app_id      bigint;
        super_admin_role_id bigint;
        admin_role_id       bigint;
        super_admin_usr_id  bigint;
    begin
        -- cfg
        insert into cfg(v)
        values ('{
          "refresh_token_dur_seconds": 1296000,
          "access_token_dur_seconds": 600
        }');

        -- app
        insert into app(name, perm_url, is_account_app)
        values ('Account', '', true)
        returning id
            into account_app_id;

        -- perms
        insert into perm(code, is_all, app_id, dsc, is_system)
        values ('acc:*', true, account_app_id, 'All permissions', true)
             , ('acc:m_app', false, account_app_id, 'Modify applications', true)
             , ('acc:m_perm', false, account_app_id, 'Modify permissions', true)
             , ('acc:m_role', false, account_app_id, 'Modify roles', true)
             , ('acc:m_usr', false, account_app_id, 'Modify users', true);

        -- Super Admin role
        insert into role(code, name, is_system)
        values ('acc:super_admin', 'Account:SuperAdmin', true)
        returning id
            into super_admin_role_id;

        -- Admin role
        insert into role(code, name, is_system)
        values ('acc:admin', 'Account:Admin', true)
        returning id
            into admin_role_id;

        -- SuperAdmin role_perm
        insert into role_perm(role_id, perm_id)
        values (super_admin_role_id, (select id from perm where app_id = account_app_id and is_all));

        -- Admin role_perm
        insert into role_perm(role_id, perm_id)
        select admin_role_id, id
        from perm
        where app_id = account_app_id
          and code in ('acc:m_role', 'acc:m_usr');

        -- SuperAdmin user
        insert into usr(phone, name)
        values ('70000000000', 'SAdmin')
        returning id
            into super_admin_usr_id;

        -- SuperAdmin usr_role
        insert into usr_role(usr_id, role_id)
        values (super_admin_usr_id, super_admin_role_id);
    end ;
$$;
