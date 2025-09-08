package entity_test

import (
	"testing"
	"time"

	"github.com/portaria-keys/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewReservation(t *testing.T) {
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	dueAt := time.Now().Add(24 * time.Hour)

	reservation := entity.NewReservation(keyID, userID, dueAt)

	assert.NotNil(t, reservation)
	assert.False(t, reservation.ID.IsZero())
	assert.Equal(t, keyID, reservation.KeyID)
	assert.Equal(t, userID, reservation.UserID)
	assert.Equal(t, dueAt.Format(time.RFC3339), reservation.DueAt.Format(time.RFC3339))
	assert.Equal(t, entity.ReservationStatusActive, reservation.Status)
	assert.Nil(t, reservation.ReturnedAt)
	assert.WithinDuration(t, time.Now(), reservation.ReservedAt, time.Second)
	assert.WithinDuration(t, time.Now(), reservation.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), reservation.UpdatedAt, time.Second)
}

func TestReservationValidateForCreation(t *testing.T) {
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// Valid Reservation
	reservation := entity.NewReservation(keyID, userID, time.Now().Add(time.Hour))
	assert.Nil(t, reservation.ValidateForCreation())

	// Invalid Reservation - DueAt in the past
	reservation = entity.NewReservation(keyID, userID, time.Now().Add(-time.Hour))
	assert.NotNil(t, reservation.ValidateForCreation())

	// Invalid Reservation - Missing KeyID
	reservation = entity.NewReservation(primitive.NilObjectID, userID, time.Now().Add(time.Hour))
	assert.NotNil(t, reservation.ValidateForCreation())

	// Invalid Reservation - Missing UserID
	reservation = entity.NewReservation(keyID, primitive.NilObjectID, time.Now().Add(time.Hour))
	assert.NotNil(t, reservation.ValidateForCreation())
}

func TestReservationValidateReturnTime(t *testing.T) {
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	reservation := entity.NewReservation(keyID, userID, time.Now().Add(time.Hour))

	// Valid return time
	returnedAt := time.Now().Add(30 * time.Minute)
	reservation.ReturnedAt = &returnedAt
	assert.Nil(t, reservation.ValidateReturnTime())

	// Invalid return time - before ReservedAt
	returnedAt = time.Now().Add(-30 * time.Minute)
	reservation.ReturnedAt = &returnedAt
	assert.NotNil(t, reservation.ValidateReturnTime())
}

func TestReservationIsOverdue(t *testing.T) {
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// Not overdue
	reservation := entity.NewReservation(keyID, userID, time.Now().Add(time.Hour))
	assert.False(t, reservation.IsOverdue())

	// Overdue
	reservation = entity.NewReservation(keyID, userID, time.Now().Add(-time.Hour))
	assert.True(t, reservation.IsOverdue())

	// Not overdue (returned)
	returnedAt := time.Now()
	reservation.ReturnedAt = &returnedAt
	reservation.Status = entity.ReservationStatusReturned
	assert.False(t, reservation.IsOverdue())
}

func TestReservationCalculateOverdueTime(t *testing.T) {
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// Not overdue
	reservation := entity.NewReservation(keyID, userID, time.Now().Add(time.Hour))
	assert.Equal(t, time.Duration(0), reservation.CalculateOverdueTime())

	// Overdue
	dueAt := time.Now().Add(-2 * time.Hour)
	reservation = entity.NewReservation(keyID, userID, dueAt)
	// The exact duration might vary slightly due to test execution time, so check if it's close
	assert.InDelta(t, 2*time.Hour, reservation.CalculateOverdueTime(), float64(time.Minute))
}

func TestReservationCanBeExtended(t *testing.T) {
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// Can be extended (active)
	reservation := entity.NewReservation(keyID, userID, time.Now().Add(time.Hour))
	reservation.Status = entity.ReservationStatusActive
	assert.True(t, reservation.CanBeExtended())

	// Cannot be extended (returned)
	reservation.Status = entity.ReservationStatusReturned
	assert.False(t, reservation.CanBeExtended())

	// Cannot be extended (overdue)
	reservation.Status = entity.ReservationStatusOverdue
	assert.False(t, reservation.CanBeExtended())
}

func TestReservationMarkAsReturned(t *testing.T) {
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	reservation := entity.NewReservation(keyID, userID, time.Now().Add(time.Hour))

	assert.Nil(t, reservation.MarkAsReturned())
	assert.Equal(t, entity.ReservationStatusReturned, reservation.Status)
	assert.NotNil(t, reservation.ReturnedAt)
	assert.WithinDuration(t, time.Now(), *reservation.ReturnedAt, time.Second)

	// Cannot mark as returned if already returned
	err := reservation.MarkAsReturned()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "reservation is not active")
}

func TestReservationMarkAsOverdue(t *testing.T) {
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	reservation := entity.NewReservation(keyID, userID, time.Now().Add(time.Hour))

	assert.Nil(t, reservation.MarkAsOverdue())
	assert.Equal(t, entity.ReservationStatusOverdue, reservation.Status)

	// Cannot mark as overdue if already overdue
	err := reservation.MarkAsOverdue()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "reservation is not active")
}
