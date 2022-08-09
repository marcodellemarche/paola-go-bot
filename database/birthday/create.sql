CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.birthday (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	"name" character varying UNIQUE NOT NULL,
    date timestamp with time zone NOT NULL,
    user_id bigint NOT NULL,
    CONSTRAINT "FK_birthday_user" FOREIGN KEY (user_id)
        REFERENCES public.user (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT name_user_id_unique_constraint UNIQUE (name, user_id)
);
