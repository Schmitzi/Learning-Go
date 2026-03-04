# 02-jwt-golang

The purpose of this exercise is to learn how `JWT` works in `GO`.

It simulates a basic authentication loop, denying users to certain endpoints when no valid `token` is provided.

### ***WARNING!***
There are some hardcoded values in this project. This is absolutely not good practise and I do not endorse ever doing this, besides within the scope of learning. If you ever do this, *NEVER* use real credentials. This is why the password is so.

## Usage

- *Welcome*: Call a Welcome Message

This is just the root. this endpoint does the following
```bash
❯ curl http://localhost:4000/

# Output
# {"status":"Success","message":"Welcome to Golang with JWT authentication"}
```

- *Login*: Login to server and create JWT token
```bash
curl -X POST http://localhost:4000/login \
  -H "Content-Type: application/json" \
  -d '{"username":"schmitzi","password":"password"}'

# Output
# {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzI3MDUxOTAsImlhdCI6MTc3MjYxODc5MCwic3ViIjoic2NobWl0emkifQ.KyC-5Yfb4gKSvrm7UomOEtvVnxg2eNHfXb32HGknqdQ"}
```

You can test the JWT token at [JWT.io](http://jwt.io)

- *Secure*: With the generated JWT token, we can add it to the headers and the access the /secure endpoint

This happens when you try to access this endpoint without a valid token:
```bash
❯ curl http://localhost:4000/secure

# Output
# {"status":"Failure","message":"You are not authorized to view this page"}
```

Then when providing a valid JWT token
```bash
curl http://localhost:4000/secure \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzI3MDUxOTAsImlhdCI6MTc3MjYxODc5MCwic3ViIjoic2NobWl0emkifQ.KyC-5Yfb4gKSvrm7UomOEtvVnxg2eNHfXb32HGknqdQ"

# Output
# {"status":"Success","message":"Congrats schmitzi and Welcome to the Secure page!"}
```
