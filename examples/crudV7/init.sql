CREATE DATABASE specgo2v7;

\c specgo2v7;

CREATE TABLE IF NOT EXISTS public.user (
    id serial PRIMARY KEY,
    login VARCHAR (50) NOT NULL,
    password VARCHAR (50) NOT NULL,
    phone VARCHAR (50),
    country VARCHAR (50) NOT NULL DEFAULT 'Moscow'
);

INSERT INTO public.user (login, password, phone, country)
VALUES ('John', '2222', '11111111', DEFAULT),
       ('Bob', '1234', '222222222', 'St-Petersburg'),
       ('Alex', '1234', '33333333', 'Pataya');

