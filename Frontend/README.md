# Frontend do Gerenciador de Chaves (React) ⚛️

Este é o frontend do sistema de gerenciamento de chaves, uma aplicação web interativa construída com React. Ele oferece uma interface amigável para visualizar, cadastrar, emprestar e devolver chaves, comunicando-se com a API do backend.

## Tecnologias Utilizadas

-   **React**: Biblioteca JavaScript para construção de interfaces de usuário.
-   **Vite**: Ferramenta de build rápida para projetos web modernos.
-   **JavaScript/JSX**: Linguagem de programação.
-   **CSS**: Estilização dos componentes.

## Estrutura do Projeto

-   `public/`: Contém arquivos estáticos como `index.html` e assets.
-   `src/`: Código fonte da aplicação.
    -   `main.jsx`: Ponto de entrada da aplicação React.
    -   `App.jsx`: Componente principal da aplicação.
    -   `assets/`: Imagens e outros recursos estáticos.
    -   `components/`: Contém os componentes reutilizáveis da interface de usuário:
        -   `Navbar.jsx`: Barra de navegação.
        -   `Footer.jsx`: Rodapé da aplicação.
        -   `KeyFormModal.jsx`: Modal para cadastrar ou editar chaves.
        -   `BorrowKeyModal.jsx`: Modal para emprestar uma chave.
        -   `ConfirmModal.jsx`: Modal de confirmação genérico.
        -   `AlertDialog.jsx`: Modal para exibir mensagens de alerta.

## Como Rodar

Para rodar o frontend, a maneira mais recomendada é utilizar o **Docker Compose** a partir da raiz do projeto, que irá orquestrar tanto o backend quanto o frontend. Consulte o `README.md` principal na raiz do projeto para instruções completas sobre como iniciar a aplicação com Docker Compose.

### Rodando Localmente (Apenas para Desenvolvimento)

Se você deseja rodar o frontend de forma isolada para desenvolvimento, siga os passos abaixo:

1.  Certifique-se de que o **Backend** esteja rodando e acessível (geralmente em `http://localhost:8080` se rodando localmente, ou o endereço do seu contêiner/serviço).
2.  Navegue até o diretório `Frontend`:
    ```bash
    cd Frontend
    ```
3.  Instale as dependências:
    ```bash
    npm install
    # ou yarn install
    ```
4.  Inicie o servidor de desenvolvimento:
    ```bash
    npm run dev
    # ou yarn dev
    ```
    A aplicação estará disponível em `http://localhost:5173` (ou outra porta disponível).

## Interação com o Backend

O frontend se comunica com o backend através de requisições HTTP. Quando rodando via Docker Compose, um servidor Nginx atua como proxy reverso. Isso significa que todas as chamadas de API do frontend são feitas para o caminho `/api/` (ex: `/api/keys`). O Nginx, configurado no arquivo `nginx.conf` dentro deste diretório, redireciona essas requisições para o serviço de backend (`http://backend:8080`) dentro da rede Docker.

Para desenvolvimento local (sem Docker Compose), as chamadas de API podem ser direcionadas diretamente para `http://localhost:8080` ou conforme configurado na variável de ambiente `VITE_API_URL`.
