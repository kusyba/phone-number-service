FROM golang:1.25.8-alpine

WORKDIR /app

RUN apk add --no-cache git postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN sqlc generate
RUN go build -o /app/server ./cmd/server

EXPOSE 8080

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
