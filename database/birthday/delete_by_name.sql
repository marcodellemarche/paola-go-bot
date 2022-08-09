DELETE FROM public.birthday
WHERE "name" = $1
RETURNING *
