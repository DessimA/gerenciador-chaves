package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserRole define os tipos de usuário
type UserRole string

const (
	UserRoleResident UserRole = "resident"
	UserRoleAdmin    UserRole = "admin"
)

// User representa um usuário do sistema
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"-"` // Nunca retornar no JSON
	Role      UserRole           `bson:"role" json:"role" validate:"required,oneof=resident admin"`
	IsBlocked bool               `bson:"is_blocked" json:"is_blocked"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewUser cria uma nova instância de User com timestamps preenchidos
func NewUser(name, email, password string, role UserRole) *User {
	now := time.Now()
	return &User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      role,
		IsBlocked: false, // Usuário não bloqueado por padrão
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ValidateForCreation valida os campos do User para criação
func (u *User) ValidateForCreation() error {
	return validate.Struct(u)
}

// ValidateForLogin valida os campos necessários para login
func (u *User) ValidateForLogin() error {
	// Apenas email e password são necessários para login
	return validate.StructPartial(u, "Email", "Password")
}

// HashPassword gera o hash da senha do usuário
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifica se a senha fornecida corresponde ao hash armazenado
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// CanMakeReservation verifica se o usuário pode fazer uma reserva
func (u *User) CanMakeReservation() bool {
	return !u.IsBlocked
}

// IsAdmin verifica se o usuário tem a role de administrador
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}