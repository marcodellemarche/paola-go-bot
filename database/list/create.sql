CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.list (
	user_id bigint NOT NULL,
	subscriber_id bigint NOT NULL,
    CONSTRAINT "FK_list_user"
        FOREIGN KEY (user_id)
            REFERENCES public.user (id),
    CONSTRAINT "FK_list_subscriber"
        FOREIGN KEY (subscriber_id)
            REFERENCES public.user (id),
    CONSTRAINT list_subscriber_id_user_id_unique_constraint UNIQUE (subscriber_id, user_id)
);
