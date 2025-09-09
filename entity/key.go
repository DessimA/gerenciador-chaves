package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

// Key representa uma chave física do prédio
type Key struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Description string             `bson:"description" json:"description" validate:"max=500"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewKey cria uma nova instância de Key com timestamps preenchidos
func NewKey(name, description string) *Key {
	now := time.Now()
	return &Key{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		IsActive:    true, // Chave ativa por padrão ao ser criada
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ValidateForCreation valida os campos da Key para criação
func (k *Key) ValidateForCreation() error {
	return validate.Struct(k)
}

// ValidateForUpdate valida os campos da Key para atualização
func (k *Key) ValidateForUpdate() error {
	// Para atualização, apenas Name e Description são validáveis diretamente
	// ID, CreatedAt, UpdatedAt, IsActive são gerenciados internamente ou não são editáveis pelo usuário
	return validate.StructPartial(k, "Name", "Description", "IsActive")
}

// CanBeReserved verifica se a chave pode ser reservada
func (k *Key) CanBeReserved() bool {
	return k.IsActive
}
