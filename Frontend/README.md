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

### Pré-requisitos

-   [Node.js](https://nodejs.org/en/download/) (versão 14 ou superior)
-   [npm](https://www.npmjs.com/get-npm) ou [Yarn](https://yarnpkg.com/)

### Instalação e Execução

1.  Certifique-se de que o **Backend** esteja rodando e acessível (geralmente em `http://localhost:8080`).
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

O frontend se comunica com o backend através de requisições HTTP para os endpoints da API. As configurações de proxy para o backend são definidas no arquivo `vite.config.js` para facilitar o desenvolvimento local, evitando problemas de CORS.