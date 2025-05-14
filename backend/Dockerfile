    FROM golang:1.24-alpine AS builder

    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    RUN go mod verify
    
    COPY . .
    
    RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/meu-servico .
    
    
    FROM alpine:latest
    
    RUN addgroup -S appgroup && adduser -S appuser -G appgroup
    
    WORKDIR /app
    
    COPY --from=builder /app/meu-servico /app/meu-servico
    
    USER appuser
    
    EXPOSE 8080
    ENV PORT=8080
    
    CMD ["/app/meu-servico"]