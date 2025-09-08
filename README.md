# Sistema de Gerenciamento de Chaves - Portaria

Este projeto implementa uma API REST para o gerenciamento de chaves em uma portaria de prédio, seguindo os princípios da Clean Architecture. A API permite o CRUD completo de chaves, sistema de reservas com controle de prazo, e funcionalidades administrativas como bloqueio de usuários e extensão de reservas.

## Funcionalidades

-   **Chaves:** Criar, listar, atualizar e remover chaves físicas.
-   **Usuários:** Registro e autenticação de usuários (moradores e administradores).
-   **Reservas:** Criar, listar, devolver e estender reservas de chaves.
-   **Controle de Acesso:** Autenticação via JWT e autorização baseada em roles (morador/administrador).
-   **Regras de Negócio:** Validação de prazos, bloqueio automático de usuários em atraso, e controle de reservas ativas por chave.

## Tecnologias Utilizadas

-   **Linguagem:** Go 1.21+
-   **Framework HTTP:** Gin Gonic
-   **Banco de Dados:** MongoDB
-   **Driver MongoDB:** Oficial MongoDB Go Driver
-   **Containerização:** Docker, Docker Compose
-   **Documentação API:** Swagger/OpenAPI 3.0
-   **Testes:** Testify

## Pré-requisitos

Antes de começar, certifique-se de ter as seguintes ferramentas instaladas:

-   [Go](https://golang.org/doc/install) (versão 1.21 ou superior)
-   [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)
-   [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Instalação e Execução

### 1. Clonar o Repositório

```bash
git clone https://github.com/portaria-keys/gerenciador-chaves.git
cd gerenciador-chaves
```

### 2. Configuração do Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis de ambiente:

```dotenv
DATABASE_URL=mongodb://localhost:27017
DATABASE_NAME=portaria_keys
SERVER_PORT=:8080
JWT_SECRET=your-super-secret-jwt-key # Altere para uma chave forte em produção
```

### 3. Execução com Docker Compose (Recomendado)

Esta é a forma mais fácil de subir a aplicação e o banco de dados.

```bash
docker-compose -f docker/docker-compose.yml up --build
```

Isso irá:
-   Construir a imagem Docker da aplicação Go.
-   Iniciar um contêiner MongoDB.
-   Iniciar um contêiner Mongo Express (interface web para o MongoDB, acessível em `http://localhost:8081`).
-   Iniciar o contêiner da aplicação Go, acessível em `http://localhost:8080`.

Para rodar em segundo plano:

```bash
docker-compose -f docker/docker-compose.yml up -d --build
```

Para parar os serviços:

```bash
docker-compose -f docker/docker-compose.yml down
```

### 4. Execução Local (Sem Docker para a Aplicação Go)

Certifique-se de ter um servidor MongoDB rodando localmente ou em um contêiner separado.

```bash
# Instalar dependências Go
go mod tidy

# Rodar a aplicação
go run cmd/server/main.go
```

## Documentação da API (Swagger/OpenAPI)

Após iniciar a aplicação (com Docker Compose ou localmente), a documentação interativa da API estará disponível em:

`http://localhost:8080/swagger/index.html`

### Gerar Documentação (se houver alterações nos comentários Swagger)

```bash
swag init -dir ./cmd/server
```

## Estrutura do Projeto

O projeto segue a Clean Architecture, com a seguinte estrutura de diretórios:

```
. (raiz do projeto)
├── cmd/
│   └── server/             # Ponto de entrada da aplicação (main.go)
├── internal/
│   ├── controller/         # Camada de Interface Adapters (HTTP handlers)
│   ├── entity/             # Camada de Entidades (regras de negócio empresariais)
│   ├── infrastructure/     # Camada de Infraestrutura (MongoDB, HTTP router, middlewares, config)
│   │   ├── config/
│   │   ├── database/
│   │   ├── http/
│   │   └── repository/
│   ├── repository/         # Camada de Interface Adapters (interfaces de repositório)
│   └── usecase/            # Camada de Casos de Uso (regras de negócio da aplicação)
├── docker/                 # Arquivos Docker (Dockerfile, docker-compose.yml)
├── docs/                   # Documentação gerada pelo Swagger
└── tests/                  # Testes (unitários, integração, API)
    ├── api/
    ├── integration/
    └── unit/
```

## Comandos de Desenvolvimento e Testes

-   **Instalar dependências:** `go mod tidy`
-   **Rodar todos os testes:** `go test ./...`
-   **Rodar testes unitários:** `go test ./tests/unit/...`
-   **Rodar testes de integração:** `go test ./tests/integration/...`
-   **Rodar testes de API:** `go test ./tests/api/...`
-   **Gerar mocks (se necessário):** `mockery --all --output tests/mocks` (requer `go install github.com/vektra/mockery/v2@latest`)

## Troubleshooting Comum

-   **`go: cannot find main module` ou `go: go.mod file not found`:** Execute `go mod init github.com/portaria-keys` na raiz do projeto.
-   **`go: command not found`:** Certifique-se de que o Go está instalado e o GOBIN está no seu PATH.
-   **`swag: command not found`:** Certifique-se de que o `swag` está instalado (`go install github.com/swaggo/swag/cmd/swag@latest`) e o GOBIN está no seu PATH.
-   **Problemas de conexão com MongoDB:** Verifique se o contêiner MongoDB está rodando (`docker ps`) e se a `DATABASE_URL` no seu `.env` está correta.

---