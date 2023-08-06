----------------forms----------------

create schema if not exists projects;

create table if not exists projects.list
(
    project_uuid uuid         not null,
    project_name varchar(100) not null,
    project_description varchar(150) not null,
    created_at   timestamp,
    deleted_at   timestamp,
    constraint pk_project primary key (project_uuid)
);
create unique index if not exists idx1_projects on projects.list (project_name);

---------------forms----------------

create table if not exists forms.project
(
    form_uuid uuid         not null,
    project_uuid uuid         not null,
    constraint pk_project primary key (form_uuid),
    constraint fk1_projects foreign key (project_uuid) references projects.list
);
create index if not exists idx1_project on forms.project (project_uuid);

---------------operations----------------

create or replace function operations.add_project(
    pProject_name varchar(100),
    pProject_description varchar(150)
) returns uuid
    language plpgsql as
$$
declare
    vResult uuid;
begin
    insert
    into projects.list
    (project_uuid,
     project_name,
     project_description,
     created_at,
     deleted_at)
    values (gen_random_uuid(),
            pProject_name,
            pProject_description,
            now(),
            null)
    returning
        project_uuid
        into
            vResult;

    return vResult;
end
$$;

create or replace function operations.add_form_project(
    pForm_uuid uuid,
    pProject_uuid uuid
) returns void
    language plpgsql as
$$
begin
    insert
    into forms.project
    (form_uuid,
     project_uuid)
    values (pForm_uuid,
            pProject_uuid);
end
$$;