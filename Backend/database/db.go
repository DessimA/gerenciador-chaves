// Pacote database é responsável pela configuração e inicialização do banco de dados.
package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // Driver SQLite em Go puro
)

var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados SQLite e cria a tabela de chaves se ela não existir.
func InitDB(filepath string) {
	var err error
	// Abre a conexão com o arquivo de banco de dados SQLite.
	DB, err = sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	// Cria a tabela de chaves se ela ainda não existir.
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS keys (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"apartment_number" TEXT,
		"key_type" TEXT,
		"status" TEXT,
		"borrowed_at" DATETIME,
		"returned_at" DATETIME,
		"borrower_name" TEXT
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar a tabela 'keys': %v", err)
	}

	log.Println("Banco de dados inicializado e tabela 'keys' criada com sucesso.")
}