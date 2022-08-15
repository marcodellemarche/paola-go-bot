DELETE FROM public.birthday
WHERE
    "name" = $1
    AND user_id = $2
RETURNING *
