CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.birthday (
	"name" character varying NOT NULL,
    contact_id bigint NULL,
    date timestamp with time zone NOT NULL,
    user_id bigint NOT NULL,
    CONSTRAINT "FK_birthday_contact"
        FOREIGN KEY (contact_id)
            REFERENCES public.user (id),
    CONSTRAINT "FK_birthday_user"
        FOREIGN KEY (user_id)
            REFERENCES public.user (id),
    CONSTRAINT birthday_name_user_id_unique_constraint UNIQUE (name, user_id)
);
