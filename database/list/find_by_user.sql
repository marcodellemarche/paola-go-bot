SELECT
    user_id,
    subscriber_id,
    user_name
FROM public.list
WHERE user_id = $1
