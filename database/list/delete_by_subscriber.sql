DELETE FROM public.list
WHERE
    subscriber_id = $1
    AND user_id = $2
RETURNING *
