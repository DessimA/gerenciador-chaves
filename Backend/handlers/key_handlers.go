// Pacote handlers contém os manipuladores de requisições HTTP.
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DessimA/gerenciador-chaves/database"
	"github.com/DessimA/gerenciador-chaves/models"

	"github.com/gorilla/mux"
)

// CreateKeyHandler cadastra uma nova chave no banco de dados.
func CreateKeyHandler(w http.ResponseWriter, r *http.Request) {
	var key models.Key
	// Decodifica o JSON do corpo da requisição para a struct Key.
	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	// Define o status inicial como "disponivel".
	key.Status = "disponivel"

	// Insere a nova chave no banco de dados.
	stmt, err := database.DB.Prepare("INSERT INTO keys(apartment_number, key_type, status) VALUES(?, ?, ?)")
	if err != nil {
		http.Error(w, "Erro ao preparar a query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(key.ApartmentNumber, key.KeyType, key.Status)
	if err != nil {
		http.Error(w, "Erro ao inserir a chave no banco", http.StatusInternalServerError)
		return
	}

	// Obtém o ID da chave inserida.
	id, _ := res.LastInsertId()
	key.ID = int(id)

	// Retorna a chave criada como JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(key)
}

// GetKeysHandler lista todas as chaves cadastradas.
func GetKeysHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, apartment_number, key_type, status, borrowed_at, returned_at, borrower_name FROM keys")
	if err != nil {
		http.Error(w, "Erro ao buscar as chaves", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var keys []models.Key
	// Itera sobre os resultados da consulta.
	for rows.Next() {
		var key models.Key
		var borrowedAt, returnedAt sql.NullTime
		var borrowerName sql.NullString

		if err := rows.Scan(&key.ID, &key.ApartmentNumber, &key.KeyType, &key.Status, &borrowedAt, &returnedAt, &borrowerName); err != nil {
			http.Error(w, "Erro ao escanear a linha", http.StatusInternalServerError)
			return
		}

		// Trata valores nulos do banco.
		if borrowedAt.Valid {
			key.BorrowedAt = &borrowedAt.Time
		}
		if returnedAt.Valid {
			key.ReturnedAt = &returnedAt.Time
		}
		if borrowerName.Valid {
			key.BorrowerName = &borrowerName.String
		}

		keys = append(keys, key)
	}

	// Retorna a lista de chaves como JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

// GetKeyHandler consulta uma chave específica pelo ID.
func GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var key models.Key
	var borrowedAt, returnedAt sql.NullTime
	var borrowerName sql.NullString

	// Consulta a chave pelo ID.
	err = database.DB.QueryRow("SELECT id, apartment_number, key_type, status, borrowed_at, returned_at, borrower_name FROM keys WHERE id = ?", id).Scan(&key.ID, &key.ApartmentNumber, &key.KeyType, &key.Status, &borrowedAt, &returnedAt, &borrowerName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Chave não encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao buscar a chave", http.StatusInternalServerError)
		}
		return
	}

	if borrowedAt.Valid {
		key.BorrowedAt = &borrowedAt.Time
	}
	if returnedAt.Valid {
		key.ReturnedAt = &returnedAt.Time
	}
	if borrowerName.Valid {
		key.BorrowerName = &borrowerName.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// BorrowKeyHandler marca uma chave como "emprestada".
func BorrowKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var payload struct {
		BorrowerName string `json:"borrower_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	if payload.BorrowerName == "" {
		http.Error(w, "O nome do retirante é obrigatório", http.StatusBadRequest)
		return
	}

	// Atualiza o status da chave para "emprestada" e registra a data/hora e o nome.
	stmt, err := database.DB.Prepare("UPDATE keys SET status = 'emprestada', borrowed_at = ?, borrower_name = ?, returned_at = NULL WHERE id = ? AND status = 'disponivel'")
	if err != nil {
		http.Error(w, "Erro ao preparar a query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(time.Now(), payload.BorrowerName, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar a chave", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Chave não está disponível para empréstimo ou não foi encontrada", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Chave %d emprestada para %s", id, payload.BorrowerName)
}

// ReturnKeyHandler marca uma chave como "disponivel".
func ReturnKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Atualiza o status da chave para "disponivel" e registra a data/hora da devolução.
	stmt, err := database.DB.Prepare("UPDATE keys SET status = 'disponivel', returned_at = ? WHERE id = ? AND status = 'emprestada'")
	if err != nil {
		http.Error(w, "Erro ao preparar a query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(time.Now(), id)
	if err != nil {
		http.Error(w, "Erro ao atualizar a chave", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Chave não estava emprestada ou não foi encontrada", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Chave %d devolvida.", id)
}

// DeleteKeyHandler exclui uma chave do banco de dados.
func DeleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("DELETE FROM keys WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Erro ao excluir a chave", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Chave não encontrada", http.StatusNotFound)
		return	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content para exclusão bem-sucedida
}

// UpdateKeyHandler atualiza os detalhes de uma chave existente.
func UpdateKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var key models.Key
	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	// Validação básica: garantir que os campos necessários estão presentes
	if key.ApartmentNumber == "" || key.KeyType == "" {
		http.Error(w, "Número do apartamento e tipo da chave são obrigatórios", http.StatusBadRequest)
		return
	}

	// Atualiza a chave no banco de dados.
	// Não permitimos a alteração do status ou informações de empréstimo por esta rota.
	stmt, err := database.DB.Prepare("UPDATE keys SET apartment_number = ?, key_type = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Erro ao preparar a query de atualização", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(key.ApartmentNumber, key.KeyType, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar a chave no banco", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Chave não encontrada", http.StatusNotFound)
		return
	}

	// Retorna a chave atualizada (ou uma confirmação)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Chave atualizada com sucesso"})
}
