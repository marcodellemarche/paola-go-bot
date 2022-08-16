SELECT
    id,
    "name"
FROM public.user
WHERE id = $1
