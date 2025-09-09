package entity_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	name := "Test User"
	email := "test@example.com"
	password := "password123"
	role := entity.UserRoleResident

	user := entity.NewUser(name, email, password, role)

	assert.NotNil(t, user)
	assert.False(t, user.ID.IsZero())
	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, password, user.Password) // Password is not hashed yet
	assert.Equal(t, role, user.Role)
	assert.False(t, user.IsBlocked)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
}

func TestUserValidateForCreation(t *testing.T) {
	// Valid User
	user := entity.NewUser("Valid User", "valid@example.com", "password123", entity.UserRoleResident)
	assert.Nil(t, user.ValidateForCreation())

	// Invalid User - Missing Name
	user = entity.NewUser("", "valid@example.com", "password123", entity.UserRoleResident)
	assert.NotNil(t, user.ValidateForCreation())

	// Invalid User - Invalid Email
	user = entity.NewUser("Valid User", "invalid-email", "password123", entity.UserRoleResident)
	assert.NotNil(t, user.ValidateForCreation())

	// Invalid User - Password too short
	user = entity.NewUser("Valid User", "valid@example.com", "short", entity.UserRoleResident)
	assert.NotNil(t, user.ValidateForCreation())

	// Invalid User - Invalid Role
	user = entity.NewUser("Valid User", "valid@example.com", "password123", "invalid_role")
	assert.NotNil(t, user.ValidateForCreation())
}

func TestUserValidateForLogin(t *testing.T) {
	// Valid Login
	user := entity.NewUser("Valid User", "valid@example.com", "password123", entity.UserRoleResident)
	assert.Nil(t, user.ValidateForLogin())

	// Invalid Login - Missing Email
	user = entity.NewUser("Valid User", "", "password123", entity.UserRoleResident)
	assert.NotNil(t, user.ValidateForLogin())

	// Invalid Login - Missing Password
	user = entity.NewUser("Valid User", "valid@example.com", "", entity.UserRoleResident)
	assert.NotNil(t, user.ValidateForLogin())
}

func TestUserHashPasswordAndCheckPassword(t *testing.T) {
	user := entity.NewUser("Test User", "test@example.com", "password123", entity.UserRoleResident)

	err := user.HashPassword()
	assert.Nil(t, err)
	assert.NotEqual(t, "password123", user.Password) // Password should be hashed

	assert.True(t, user.CheckPassword("password123"))
	assert.False(t, user.CheckPassword("wrong_password"))
}

func TestUserCanMakeReservation(t *testing.T) {
	user := entity.NewUser("Test User", "test@example.com", "password123", entity.UserRoleResident)
	assert.True(t, user.CanMakeReservation())

	user.IsBlocked = true
	assert.False(t, user.CanMakeReservation())
}

func TestUserIsAdmin(t *testing.T) {
	user := entity.NewUser("Test User", "test@example.com", "password123", entity.UserRoleResident)
	assert.False(t, user.IsAdmin())

	user.Role = entity.UserRoleAdmin
	assert.True(t, user.IsAdmin())
}
