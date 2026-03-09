CREATE TABLE IF NOT EXISTS books (
    id       SERIAL PRIMARY KEY,
    title    VARCHAR(255),
    author   VARCHAR(255),
    price    NUMERIC(10,2),
    imageurl VARCHAR(255)
);
