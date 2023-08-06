create function get_form(pform_uuid uuid)
    returns TABLE(data_uuid uuid, field_name character varying, field_type character varying, field_data character varying)
    language plpgsql
as
$$
begin
    return query
    select fd.data_uuid,
           fd.field_name,
           fd.field_type,
           fd.field_data
    from forms.data as fd
    where form_uuid = pForm_uuid;
end
$$;

create function get_attachment(pform_uuid uuid)
    returns TABLE(file_uuid uuid, file_name character varying, file_type character varying)
    language plpgsql
as
$$
begin
    return query
    select at.file_uuid,
           at.file_name,
           at.file_type
    from forms.attachments as at
    where form_uuid = pForm_uuid;
end
$$;

create function get_form_list(plimit bigint, poffset bigint)
    returns TABLE(form_uuid uuid, form_id smallint, created_at timestamp without time zone)
    language plpgsql
as
$$
begin
    return query
    select fl.form_uuid,
           fl.form_id,
           fl.created_at
    from forms.list as fl
    limit pLimit offset pOffset;

end
$$;