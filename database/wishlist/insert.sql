INSERT INTO public.wishlist
(
    user_id,
    "name",
    link,
    buyer_id
)
VALUES
(
    $1,
    $2,
    $3,
    $4
)
RETURNING *
