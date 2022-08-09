DELETE FROM public.user
WHERE "id" = $1
RETURNING *
