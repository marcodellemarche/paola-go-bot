CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.wishlist (
	user_id bigint NOT NULL,
    "name" character varying NOT NULL,
    link character varying NULL,
	buyer_id bigint NULL,
    CONSTRAINT "FK_wishlist_user"
        FOREIGN KEY (user_id)
            REFERENCES public.user (id),
    CONSTRAINT "FK_wishlist_buyer"
        FOREIGN KEY (buyer_id)
            REFERENCES public.user (id),
    CONSTRAINT wishlist_name_user_id_unique_constraint UNIQUE (name, user_id)
);
