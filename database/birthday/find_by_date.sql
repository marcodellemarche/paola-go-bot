SELECT "name", contact_id, date, user_id
FROM birthday
WHERE
	(SELECT DATE_PART('day', "date")) = $1
	AND (SELECT DATE_PART('month', "date")) = $2;
	