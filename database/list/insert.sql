INSERT INTO public.list
(
    user_id,
    subscriber_id,
    user_name
)
VALUES
(
    $1,
    $2,
    $3
)
ON CONFLICT DO NOTHING
RETURNING *
