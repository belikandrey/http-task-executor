## Task executor Rest API

### Full list what has been used:
* [go-chi](https://github.com/go-chi/chi) Router for HTTP service
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [goose](https://github.com/pressly/goose) - DB migrations
* [swag](https://github.com/swaggo/swag) - Swagger
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [gomock](https://github.com/golang/mock) - Mocking framework
* [cleanenv](https://github.com/ilyakaznacheev/cleanenv) - For config



#### Generate swagger docs swag init -g cmd/api/main.go
#### Swagger available by default at => http://localhost:8081/swagger-ui

#### Run with flag --config=./config/local.yaml(prod.yaml) or with env variable CONFIG_PATH (default => http://localhost:8081)