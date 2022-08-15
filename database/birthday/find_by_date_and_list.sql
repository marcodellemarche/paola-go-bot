SELECT birthday."name", birthday.contact_id, birthday.date, list.subscriber_id AS user_id
FROM public.list
INNER JOIN birthday
    ON list.user_id = birthday.user_id
WHERE
    (SELECT DATE_PART('day', birthday.date)) = $1
	AND (SELECT DATE_PART('month', birthday.date)) = $2
    AND list.user_id = $3
    AND (CAST($4 AS BIGINT) IS NULL OR list.subscriber_id = CAST($4 AS BIGINT));
