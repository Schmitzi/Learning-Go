# 03-dockerize-go

This project builds on `01-bookstore-api` by containerising the application using Docker. The Go REST API and PostgreSQL database are orchestrated with docker-compose, allowing the entire stack to be spun up with a single command.

The project also introduces two important production concepts: TLS encryption for secure HTTPS communication, and Swagger/OpenAPI documentation for interactive API exploration and testing..

## Usage

PostgreSQL is still used for the database connection, but now its managed in our `docker.compose.yml` file.
```docker
postgresql:
    image: postgres:15.17
    container_name: postgresql
    env_file:
      - ./.env
    networks:
        - bookstore-api
    volumes:
      - bookstore:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped
```

This will create a database container with a `database` named `bookstore`. Then using our `init.sql` file to create the `table` named `books`.
```sql
CREATE TABLE IF NOT EXISTS books (
    id       SERIAL PRIMARY KEY,
    title    VARCHAR(255),
    author   VARCHAR(255),
    price    NUMERIC(10,2),
    imageurl VARCHAR(255)
);
```

You will also need a `.env` file containing:
```text
POSTGRES_DB: bookstore
POSTGRES_USER: user
POSTGRES_PASSWORD: password

DATABASE_URL: postgres://user:password@postgresql/bookstore?sslmode=disable
PORT: 8080
```

You will need to generate `keys.pem & cert.pem` for the `TLS` connection. The command to generate these is 
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout key.pem -out cert.pem \
    -subj "/C=NZ/ST=Auckland/L=Auckland/O=bookstore/CN=localhost"
```

To launch the server, just build the container and restart
```bash
docker compose up --build
```

This starts the server listening on the 8080 port. If you want to change it then it needs to be changed in the `docker-compose.yml`

As before, the previous HTTP methods exist

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
