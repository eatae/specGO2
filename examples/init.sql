CREATE DATABASE specgo;

\c specgo;

CREATE TABLE IF NOT EXISTS articles (
     id serial PRIMARY KEY,
     title VARCHAR (50) NOT NULL,
     author VARCHAR (50) NOT NULL,
     content VARCHAR (50) NOT NULL
);

INSERT INTO articles (title, author, content)
    VALUES ('First title', 'Gnom', 'Some long content.'),
        ('Second title', 'Claus', 'Very long content.');


CREATE TABLE IF NOT EXISTS users (
     id serial PRIMARY KEY,
     login VARCHAR (50) NOT NULL,
     password VARCHAR (50) NOT NULL
);

INSERT INTO users (login, password)
    VALUES ('John', '2222'),
        ('Bob', '1234');