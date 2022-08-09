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
RETURNING *
