INSERT INTO public.birthday
(
    "name",
    date,
    user_id
)
VALUES
(
    $1,
    $2,
    $3
)
RETURNING *
