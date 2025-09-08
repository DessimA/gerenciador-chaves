package entity_test

import (
	"strings"
	"testing"
	"time"

	"github.com/portaria-keys/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewKey(t *testing.T) {
	name := "Test Key"
	description := "A key for testing purposes"

	key := entity.NewKey(name, description)

	assert.NotNil(t, key)
	assert.False(t, key.ID.IsZero())
	assert.Equal(t, name, key.Name)
	assert.Equal(t, description, key.Description)
	assert.True(t, key.IsActive)
	assert.WithinDuration(t, time.Now(), key.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), key.UpdatedAt, time.Second)
}

func TestKeyValidateForCreation(t *testing.T) {
	// Valid Key
	key := entity.NewKey("Valid Key", "Description")
	assert.Nil(t, key.ValidateForCreation())

	// Invalid Key - Missing Name
	key = entity.NewKey("", "Description")
	assert.NotNil(t, key.ValidateForCreation())

	// Invalid Key - Name too short
	key = entity.NewKey("a", "Description")
	assert.NotNil(t, key.ValidateForCreation())

	// Invalid Key - Name too long
	key = entity.NewKey(strings.Repeat("a", 101), "Description")
	assert.NotNil(t, key.ValidateForCreation())

	// Invalid Key - Description too long
	key = entity.NewKey("Valid Key", strings.Repeat("a", 501))
	assert.NotNil(t, key.ValidateForCreation())
}

func TestKeyValidateForUpdate(t *testing.T) {
	// Valid Update
	key := entity.NewKey("Valid Key", "Description")
	key.Name = "Updated Name"
	key.Description = "Updated Description"
	key.IsActive = false
	assert.Nil(t, key.ValidateForUpdate())

	// Invalid Update - Name too short
	key.Name = "a"
	assert.NotNil(t, key.ValidateForUpdate())
}

func TestKeyCanBeReserved(t *testing.T) {
	key := entity.NewKey("Test Key", "Description")
	assert.True(t, key.CanBeReserved())

	key.IsActive = false
	assert.False(t, key.CanBeReserved())
}
