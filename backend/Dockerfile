FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server cmd/webserver/main.go

# Build goose for migrations
RUN go build -o /goose -tags='no_clickhouse no_mssql no_mysql no_sqlite3 no_vertica no_ydb' \
    github.com/pressly/goose/v3/cmd/goose

# Runtime
FROM alpine:3.21

WORKDIR /

COPY --from=builder /server /server
COPY --from=builder /goose /goose
COPY --from=builder /app/db/migrations /migrations

RUN adduser -D nonroot nonroot
RUN chown nonroot:nonroot /server /goose
USER nonroot:nonroot

EXPOSE 8008

ENTRYPOINT ["/server"]
