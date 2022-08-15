SELECT user_id, subscriber_id
FROM public.list
WHERE user_id = $1
