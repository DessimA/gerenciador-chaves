# Gerenciador de Chaves 🔑

Bem-vindo ao **Gerenciador de Chaves**! Este é um projeto didático e divertido para você aprender como construir uma aplicação completa, desde o backend robusto em Go até um frontend interativo em React.

Já imaginou ter um sistema para controlar as chaves de um prédio, sabendo quem pegou qual chave, quando foi emprestada e quando foi devolvida? É exatamente isso que este projeto faz! Ele simula um sistema simples de gerenciamento de chaves, perfeito para quem está dando os primeiros passos em desenvolvimento full-stack.

## O que você vai encontrar aqui?

-   **Backend em Go**: Uma API REST simples e eficiente, construída com Go, que lida com toda a lógica de negócio e persistência de dados.
-   **Frontend em React**: Uma interface de usuário moderna e responsiva, desenvolvida com React, para você interagir com o sistema de forma intuitiva.
-   **Dockerização**: Tudo empacotado em containers Docker para facilitar a execução e o deploy.

## Como Rodar a Aplicação Completa com Docker Compose

Para facilitar a execução de toda a aplicação (Backend e Frontend) de uma só vez, você pode usar o Docker Compose.

### Pré-requisitos

-   [Docker Desktop](https://www.docker.com/products/docker-desktop) (inclui Docker Engine e Docker Compose)

### Passos para Execução

1.  **Navegue até a raiz do projeto**:
    Abra seu terminal ou prompt de comando e vá para o diretório principal do projeto `gerenciador-chaves` (onde o arquivo `docker-compose.yml` está localizado).
    ```bash
    cd D:\Github\gerenciador-chaves
    ```

2.  **Construa as imagens Docker**:
    Este comando irá construir as imagens para o backend (Go) e para o frontend (React/Nginx) com base nos `Dockerfiles` e no `docker-compose.yml`.
    ```bash
    docker-compose build
    ```
    *Aguarde a conclusão do processo. Isso pode levar alguns minutos na primeira vez.*

3.  **Inicie os contêineres**:
    Após a construção das imagens, este comando irá iniciar os serviços `backend` e `frontend` em segundo plano.
    ```bash
    docker-compose up -d
    ```
    *O `-d` no final significa "detached mode", ou seja, os contêineres rodarão em segundo plano.*

4.  **Acesse a Aplicação**:
    -   **Frontend**: Abra seu navegador e acesse `http://localhost` (ou `http://localhost:80`).
    -   **Backend (API)**: A API estará disponível em `http://localhost:8080`. Você pode testá-la usando ferramentas como `curl` ou Postman.

### Parando a Aplicação

Para parar e remover os contêineres (mas manter as imagens e o volume do banco de dados), execute na raiz do projeto:

```bash
docker-compose down
```

### Persistência de Dados

O banco de dados SQLite (`keys.db`) do backend é persistido através de um volume Docker. Isso significa que seus dados não serão perdidos quando você parar ou reiniciar os contêineres.
