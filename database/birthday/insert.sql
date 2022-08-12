INSERT INTO public.birthday
(
    "name",
    contact_id,
    date,
    user_id
)
VALUES
(
    $1,
    $2,
    $3,
    $4
)
RETURNING *
