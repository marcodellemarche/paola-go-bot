SELECT id, "name", contact_id, date, user_id
FROM public.birthday
WHERE user_id = $1
