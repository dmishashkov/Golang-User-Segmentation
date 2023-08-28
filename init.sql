CREATE DATABASE avito2023;


\c avito2023;

CREATE TABLE public.slugs (
    slug_id  bigserial PRIMARY KEY,
    slug_name TEXT
);


CREATE TABLE public.users (
    user_id bigint PRIMARY KEY
);

CREATE TABLE public.slugs_users (
    user_id bigint,
    slug_id bigint
);
