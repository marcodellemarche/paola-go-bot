UPDATE public.wishlist
SET
    buyer_id = $3
WHERE
    "name" = $1
    AND user_id = $2
