# XM Golang Exercise

In order to run the program, it is necessary to first prepare the database.

This can be done with the commands
```bash
docker run -d -p 5432:5432 --name xmdb -e POSTGRES_PASSWORD=xm postgres
export PGPASSWORD='xm'
psql -h localhost -p 5432 -U postgres -f init.sql
```

It is then necessary to download the dependencies
```bash
go get github.com/gin-gonic/gin
go get github.com/google/uuid
go get github.com/jackc/pgx/stdlib
```

Once the dependencies are downloaded, the program can be run with
```bash
go run main.go
```

It will listen on port `8080`.

The program will run in debug mode by default, to run in production mode, set the `GIN_MODE` environment variable to `release`.

In order to authenticate a request with curl, specify the `--user "admin:xm"` flag.

For example
```bash
curl -X POST localhost:8080/ -d "{
  \"name\": \"test\",
  \"amount_of_employees\": 1,
  \"registered\": true,
  \"company_type\": \"Cooperative\" }" --user "admin:xm"
```
