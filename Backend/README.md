# Backend do Gerenciador de Chaves (Go) 🚀

Este é o coração do sistema de gerenciamento de chaves, uma API RESTful construída em Go. Ele é responsável por toda a lógica de negócio, manipulação de dados e interação com o banco de dados.

## Tecnologias Utilizadas

-   **Go**: Linguagem de programação principal.
-   **Gorilla Mux**: Roteador HTTP para lidar com as requisições da API.
-   **SQLite**: Banco de dados leve e simples, ideal para este projeto didático.

## Estrutura do Projeto

-   `main.go`: Ponto de entrada da aplicação, onde o servidor HTTP é configurado e as rotas são definidas.
-   `models/key.go`: Define a estrutura de dados para uma chave (`Key`) e seus métodos associados.
-   `handlers/key_handlers.go`: Contém as funções (handlers) que processam as requisições HTTP para cada endpoint da API, interagindo com o banco de dados.
-   `database/db.go`: Responsável pela inicialização do banco de dados SQLite e pela criação da tabela `keys` se ela não existir.
-   `Dockerfile`: Arquivo para construir a imagem Docker do backend.

## Modelagem de Dados (SQLite)

O banco de dados utiliza uma única tabela chamada `keys` com a seguinte estrutura:

| Campo            | Tipo      | Descrição                               |
| :--------------- | :-------- | :-------------------------------------- |
| `id`             | INTEGER   | Chave primária, auto-incremento         |
| `apartment_number` | TEXT      | Número do apartamento associado à chave |
| `key_type`       | TEXT      | Tipo da chave (ex: "apartamento", "garagem", "deposito") |
| `status`         | TEXT      | Status da chave ("disponivel", "emprestada") |
| `borrowed_at`    | DATETIME  | Timestamp de quando a chave foi emprestada (pode ser NULL) |
| `returned_at`    | DATETIME  | Timestamp de quando a chave foi devolvida (pode ser NULL) |
| `borrower_name`  | TEXT      | Nome de quem pegou a chave (pode ser NULL) |

## Endpoints da API

A API expõe os seguintes endpoints:

### `GET /keys`

-   **Descrição**: Lista todas as chaves cadastradas no sistema.
-   **Resposta**: Um array de objetos `Key`.

### `POST /keys`

-   **Descrição**: Cadastra uma nova chave.
-   **Corpo da Requisição (JSON)**:
    ```json
    {
        "apartment_number": "string",
        "key_type": "string"
    }
    ```
-   **Resposta**: O objeto `Key` da chave recém-criada.

### `GET /keys/{id}`

-   **Descrição**: Consulta os detalhes de uma chave específica pelo seu `id`.
-   **Parâmetros de URL**:
    -   `id`: ID da chave (inteiro).
-   **Resposta**: O objeto `Key` correspondente.

### `PUT /keys/{id}/borrow`

-   **Descrição**: Marca uma chave como emprestada.
-   **Parâmetros de URL**:
    -   `id`: ID da chave (inteiro).
-   **Corpo da Requisição (JSON)**:
    ```json
    {
        "borrower_name": "string"
    }
    ```
-   **Resposta**: O objeto `Key` atualizado.

### `PUT /keys/{id}/return`

-   **Descrição**: Marca uma chave como devolvida.
-   **Parâmetros de URL**:
    -   `id`: ID da chave (inteiro).
-   **Resposta**: O objeto `Key` atualizado.

## Como Rodar

### Pré-requisitos

-   [Go](https://golang.org/doc/install) (versão 1.16 ou superior)
-   [Docker](https://docs.docker.com/get-docker/) (opcional, para rodar via container)

### Rodando Localmente

1.  Navegue até o diretório `Backend`:
    ```bash
    cd Backend
    ```
2.  Baixe as dependências:
    ```bash
    go mod tidy
    ```
3.  Execute a aplicação:
    ```bash
    go run main.go
    ```
    A API estará disponível em `http://localhost:8080`.

### Rodando com Docker

1.  Navegue até o diretório `Backend`:
    ```bash
    cd Backend
    ```
2.  Construa a imagem Docker:
    ```bash
    docker build -t key-manager-backend .
    ```
3.  Execute o container, mapeando a porta 8080:
    ```bash
    docker run -p 8080:8080 -v $(pwd)/keys.db:/app/keys.db key-manager-backend
    ```
    *O volume `-v $(pwd)/keys.db:/app/keys.db` garante que o banco de dados `keys.db` seja persistido no seu diretório local, mesmo que o contêiner seja removido.*
    A API estará disponível em `http://localhost:8080`.

## Exemplos de Uso (com `curl`)

```bash
# Cadastrar uma nova chave
curl -X POST http://localhost:8080/keys -H "Content-Type: application/json" -d '{"apartment_number":"101","key_type":"apartamento"}'

# Listar todas as chaves
curl http://localhost:8080/keys

# Consultar uma chave específica (substitua {id} pelo ID real da chave)
curl http://localhost:8080/keys/{id}

# Emprestar uma chave (substitua {id} pelo ID real da chave)
curl -X PUT http://localhost:8080/keys/{id}/borrow -H "Content-Type: application/json" -d '{"borrower_name":"João Silva"}'

# Devolver uma chave (substitua {id} pelo ID real da chave)
curl -X PUT http://localhost:8080/keys/{id}/return
```