# Estágio 1: Build
# Usa a imagem oficial do Go como base.
FROM golang:1.19-alpine AS builder

# Define o diretório de trabalho dentro do container.
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências.
COPY go.mod ./
COPY go.sum ./

# Baixa as dependências.
RUN go mod download

# Copia todo o código fonte do projeto.
COPY . .

# Compila a aplicação.
# CGO_ENABLED=0 desabilita o CGO para criar um binário estático.
# -o /app/main cria o binário chamado 'main' no diretório /app.
RUN CGO_ENABLED=0 go build -o /app/main ./main.go

# Estágio 2: Run
# Usa uma imagem Alpine mínima para a execução, reduzindo o tamanho final.
FROM alpine:latest

# Define o diretório de trabalho.
WORKDIR /app

# Copia o binário compilado do estágio de build.
COPY --from=builder /app/main .

# Copia o arquivo de banco de dados (se existir localmente, para iniciar com dados).
# Se não existir, a aplicação criará um novo.
COPY keys.db .

# Expõe a porta 8080 que a aplicação usa.
EXPOSE 8080

# Define o comando para executar a aplicação quando o container iniciar.
CMD ["./main"]