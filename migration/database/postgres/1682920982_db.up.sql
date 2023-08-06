drop schema if exists forms cascade;

CREATE EXTENSION "uuid-ossp";

----------------forms----------------

create schema if not exists forms;

create table if not exists forms.list
(
    form_uuid  uuid not null,
    form_id    smallint                       not null,
    created_at timestamp,
    constraint pk_form primary key (form_uuid)
);

create index if not exists idx1_forms on forms.list (form_id);

create table if not exists forms.data
(
    data_uuid  uuid not null,
    form_uuid  uuid                           not null,
    field_name varchar(100)                   not null,
    field_type varchar(100)                   not null,
    field_data character varying              not null,
    constraint pk_data primary key (data_uuid),
    constraint fk1_forms foreign key (form_uuid) references forms.list
);

create table if not exists forms.attachments
(
    file_uuid uuid not null,
    form_uuid uuid                           not null,
    file_name varchar(150)                   not null,
    file_type varchar(100)                   not null,
    constraint pk_file primary key (file_uuid),
    constraint fk1_forms foreign key (form_uuid) references forms.list
);

create index if not exists idx1_attachments on forms.attachments (form_uuid);

---------------operations----------------

create schema if not exists operations;

create or replace function operations.add_form(
    pForm_id smallint
) returns uuid
    language plpgsql as
$$
declare
    vResult uuid;
begin
    insert
    into forms.list
    (form_uuid,
     form_id,
     created_at)
    values (gen_random_uuid(),
            pForm_id,
            now())
        returning
            form_uuid
    into
        vResult;

    return vResult;
end
$$;

create or replace function operations.add_data(
    pForm_uuid uuid,
    pField_name varchar(100),
    pField_type varchar(100),
    pField_data character varying
) returns uuid
    language plpgsql as
$$
declare
    vResult uuid;
begin
    insert
    into forms.data
    (data_uuid,
     form_uuid,
     field_name,
     field_type,
     field_data)
    values (gen_random_uuid(),
            pForm_uuid,
            pField_name,
            pField_type,
            pField_data)
        returning
            data_uuid
    into
        vResult;

    return vResult;
end
$$;

create or replace function operations.add_attachment(
    pForm_uuid uuid,
    pFile_name varchar(150),
    pFile_type varchar(100)
) returns uuid
    language plpgsql as
$$
declare
    vResult uuid;
begin
    insert
    into forms.attachments
    (file_uuid,
     form_uuid,
     file_name,
     file_type)
    values (gen_random_uuid(),
            pForm_uuid,
            pFile_name,
            pFile_type)
        returning
            file_uuid
    into
        vResult;

return vResult;
end
$$;
