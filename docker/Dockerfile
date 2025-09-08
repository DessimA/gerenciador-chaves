FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod ./go.mod
RUN go mod download && go mod verify

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

CMD ["./main"]
