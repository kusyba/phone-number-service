FROM golang:1.25.8-alpine

WORKDIR /app

RUN apk add --no-cache git postgresql-client curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Скачивание совместимой версии sqlc
RUN curl -L https://github.com/sqlc-dev/sqlc/releases/download/v1.27.0/sqlc_1.27.0_linux_amd64.tar.gz -o sqlc.tar.gz && \
    tar -xzf sqlc.tar.gz && \
    mv sqlc /usr/local/bin/ && \
    rm -f sqlc.tar.gz

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# НЕ ЗАПУСКАЕМ sqlc generate здесь!
# Просто копируем entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
