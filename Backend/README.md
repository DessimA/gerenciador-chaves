# Gerenciador de Chaves de Prédio

Este é um projeto didático em Go que implementa uma API REST simples para gerenciar o empréstimo e a devolução de chaves de um prédio.

## Funcionalidades

*   Cadastrar uma nova chave.
*   Listar todas as chaves.
*   Consultar o status de uma chave específica.
*   Emprestar uma chave.
*   Devolver uma chave.

## Tecnologias Utilizadas

*   **Linguagem:** Go
*   **Roteamento HTTP:** `gorilla/mux`
*   **Banco de Dados:** SQLite
*   **Containerização:** Docker

## Estrutura do Projeto

```
/
├── main.go               # Ponto de entrada da aplicação
├── models/
│   └── key.go            # Estrutura de dados da chave
├── handlers/
│   └── key_handlers.go   # Handlers para as rotas da API
├── database/
│   └── db.go             # Configuração do banco de dados
├── go.mod                # Dependências do Go
├── Dockerfile            # Containerização da aplicação
└── README.md             # Documentação
```

## API Endpoints

| Método | Rota                  | Descrição                               | Corpo da Requisição (Exemplo)                               |
| :----- | :-------------------- | :-------------------------------------- | :---------------------------------------------------------- |
| `POST` | `/keys`               | Cadastra uma nova chave                 | `{"apartment_number": "101", "key_type": "apartamento"}`    |
| `GET`  | `/keys`               | Lista todas as chaves                   | N/A                                                         |
| `GET`  | `/keys/{id}`          | Consulta uma chave específica           | N/A                                                         |
| `PUT`  | `/keys/{id}/borrow`   | Empresta uma chave                      | `{"borrower_name": "João Silva"}`                           |
| `PUT`  | `/keys/{id}/return`   | Devolve uma chave                       | N/A                                                         |

## Como Executar

### Usando Go (Localmente)

1.  **Clone o repositório:**
    ```bash
    git clone https://github.com/seu-usuario/gerenciador-chaves.git
    cd gerenciador-chaves
    ```

2.  **Instale as dependências:**
    ```bash
    go mod tidy
    ```

3.  **Execute a aplicação:**
    ```bash
    go run main.go
    ```
    O servidor estará rodando em `http://localhost:8080`.

### Usando Docker

1.  **Construa a imagem Docker:**
    ```bash
    docker build -t gerenciador-chaves .
    ```

2.  **Execute o container:**
    ```bash
    docker run -p 8080:8080 gerenciador-chaves
    ```
    O servidor estará acessível em `http://localhost:8080`.

## Exemplos de Uso com `curl`

*   **Cadastrar uma nova chave:**
    ```bash
    curl -X POST http://localhost:8080/keys -H "Content-Type: application/json" -d '{"apartment_number": "205", "key_type": "garagem"}'
    ```

*   **Listar todas as chaves:**
    ```bash
    curl http://localhost:8080/keys
    ```

*   **Emprestar a chave com ID 1:**
    ```bash
    curl -X PUT http://localhost:8080/keys/1/borrow -H "Content-Type: application/json" -d '{"borrower_name": "Maria Souza"}'
    ```

*   **Devolver a chave com ID 1:**
    ```bash
    curl -X PUT http://localhost:8080/keys/1/return
    ```
