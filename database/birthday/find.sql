SELECT
	"name",
	contact_id,
	date,
	user_id,
	NULL AS list_id,
	NULL AS list_name
FROM birthday
WHERE
	(CAST($1 AS SMALLINT) IS NULL OR (SELECT DATE_PART('day', date)) = CAST($1 AS SMALLINT))
    AND (CAST($2 AS SMALLINT) IS NULL OR (SELECT DATE_PART('month', date)) = CAST($2 AS SMALLINT))
	AND (CAST($3 AS BIGINT) IS NULL OR user_id = CAST($3 AS BIGINT));
	