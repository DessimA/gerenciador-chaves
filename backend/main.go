package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Requisição recebida em /hello")	
	fmt.Fprint(w, "Olá mundo testando hot reload, funcionou!?")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Requisição recebida em /health")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"OK"}`)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/health", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
		log.Printf("PORT não definida, usando padrão DENTRO do container: %s", port)
	}

	addr := ":" + port
	log.Printf("Servidor escutando na porta INTERNA %s...", port)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}