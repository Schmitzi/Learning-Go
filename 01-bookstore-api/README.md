# 01-bookstore-api

Lets be honest, I got carried away with this one. The difficulty curve between the first example and this is steep.
But I felt these are necessary skills to practice.

I have converted the first exercise from static JSON to a more dynamic REST API with `Chi` for routing and `PostgreSQL` for the database

## Usage

PostgreSQL is needed for the database connection. The database can be created like so:
```sql
CREATE DATABASE bookstore;
\c bookstore

CREATE TABLE books (
    id       SERIAL PRIMARY KEY,
    title    VARCHAR(255),
    author   VARCHAR(255),
    price    NUMERIC(10,2),
    imageurl VARCHAR(255)
);
```

This will create a database called `bookstore` with a table called `books`. Then the server can connect to the database.

You will also need a `.env` file containing:
```bash
DATABASE_URL=postgres://user:password@localhost/bookstore?sslmode=disable
```

To launch the server, just build and start
```bash
go build main.go
./main
```

This starts the server listening on the 8080 port (can be changed by editing the `const` in `main.go`)

As before, the previous functions exist but are now HTTP methods

- *Add:* Add book to the library
```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"test_book","author":"Me","price":0,"imageurl":"http://test.com"}'
```

- *Get*: Get all books or book by `id`
```bash
# All books
curl http://localhost:8080/books

# By ID
curl http://localhost:8080/books/1

# Output
{"id":1,"title":"test_book","author":"Me","price":0,"imageurl":"http://test.com"}
```

- *Update*: Update book by `id`
```bash
curl -X PUT http://localhost:8080/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"test_book","author":"Me","price":1000,"imageurl":"http://test.com"}'
```

- *Delete*: Delete book by `id`
```bash
curl -X DELETE http://localhost:8080/books/1
```
