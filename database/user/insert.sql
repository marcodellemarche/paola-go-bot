INSERT INTO public.user
(
    "id",
    "name"
)
VALUES
(
    $1,
    $2
)
ON CONFLICT DO NOTHING
RETURNING *
