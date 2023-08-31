CREATE DATABASE avito2023;


\c avito2023;



CREATE TABLE public.users (
    user_id bigint PRIMARY KEY
);

CREATE TABLE public.segments (
    segment_id  bigserial PRIMARY KEY,
    segment_name TEXT
);


CREATE TABLE public.segments_users (
    user_id bigint,
    segment_id bigint,
    PRIMARY KEY (user_id, segment_id)
);


CREATE TABLE public.segments_history (
    id  bigserial PRIMARY KEY,
    user_id bigint,
    segment_id bigint,
    action_date timestamp,
    action_type TEXT
);

