DELETE FROM public.wishlist
WHERE
    name = $1
    AND user_id = $2
RETURNING *
