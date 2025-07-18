# 1: сборка бинарника

FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /todo_list_app

# 2: сборка финального образа

FROM alpine:3.22.0

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /todo_list_app /app/todo_list_app

COPY web /app/web

ENV TODO_PORT=7540 \
    TODO_DBFILE=/data/scheduler.db \
    TODO_PASSWORD=hey_Practicum!
    
EXPOSE ${TODO_PORT}

CMD ["./todo_list_app"]