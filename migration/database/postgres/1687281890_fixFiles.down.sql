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