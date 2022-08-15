INSERT INTO public.list
(
    user_id,
    subscriber_id
)
VALUES
(
    $1,
    $2
)
ON CONFLICT DO NOTHING
RETURNING *
