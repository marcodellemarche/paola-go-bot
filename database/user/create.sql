CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.user (
	id bigint NOT NULL PRIMARY KEY,
	"name" character varying UNIQUE NOT NULL
);
