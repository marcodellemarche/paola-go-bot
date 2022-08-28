SELECT
    user_id,
    "name",
    link,
    buyer_id
FROM public.wishlist
WHERE user_id = $1
