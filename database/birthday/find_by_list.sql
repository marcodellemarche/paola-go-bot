SELECT
    birthday."name",
    birthday.contact_id,
    birthday.date, 
    list.subscriber_id AS user_id,
    list.user_id AS list_id,
    list.user_name AS list_name
FROM public.list
INNER JOIN birthday
    ON list.user_id = birthday.user_id
WHERE
    (CAST($1 AS SMALLINT) IS NULL OR (SELECT DATE_PART('day', birthday.date)) = CAST($1 AS SMALLINT))
    AND (CAST($2 AS SMALLINT) IS NULL OR (SELECT DATE_PART('month', birthday.date)) = CAST($2 AS SMALLINT))
    AND (CAST($3 AS BIGINT) IS NULL OR list.user_id = CAST($3 AS BIGINT))
    AND (CAST($4 AS BIGINT) IS NULL OR list.subscriber_id = CAST($4 AS BIGINT));
