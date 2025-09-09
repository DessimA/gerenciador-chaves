// Pacote main é o ponto de entrada da aplicação.
package main

import (
	"log"
	"net/http"

	"github.com/DessimA/gerenciador-chaves/database"
	"github.com/DessimA/gerenciador-chaves/handlers"

	cors "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Inicializa o banco de dados. O arquivo do SQLite será criado no diretório atual.
	database.InitDB("./keys.db")

	// Cria um novo roteador usando gorilla/mux.
	r := mux.NewRouter()

	// Define as rotas da API.
	r.HandleFunc("/keys", handlers.GetKeysHandler).Methods("GET")
	r.HandleFunc("/keys", handlers.CreateKeyHandler).Methods("POST")
	r.HandleFunc("/keys/{id}", handlers.GetKeyHandler).Methods("GET")
	r.HandleFunc("/keys/{id}", handlers.UpdateKeyHandler).Methods("PUT") // New
	r.HandleFunc("/keys/{id}", handlers.DeleteKeyHandler).Methods("DELETE") // New
	r.HandleFunc("/keys/{id}/borrow", handlers.BorrowKeyHandler).Methods("PUT")
	r.HandleFunc("/keys/{id}/return", handlers.ReturnKeyHandler).Methods("PUT")

	// Configurações do CORS
	corsHandler := cors.CORS(
		cors.AllowedOrigins([]string{"*"}), // Em produção, restrinja para o seu domínio do frontend
		cors.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		cors.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Inicia o servidor na porta 8080 com o middleware CORS.
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(r)))
}