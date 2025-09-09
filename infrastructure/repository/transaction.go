package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrConcurrentModification indica que houve uma alteração concorrente no documento
	ErrConcurrentModification = errors.New("modificação concorrente detectada")
)

// TransactionManager gerencia transações MongoDB
type TransactionManager struct {
	client *mongo.Client
}

// NewTransactionManager cria uma nova instância do TransactionManager
func NewTransactionManager(client *mongo.Client) *TransactionManager {
	return &TransactionManager{
		client: client,
	}
}

// WithTransaction executa uma função dentro de uma transação
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	session, err := tm.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	return err
}

// IsOptimisticLockError verifica se um erro é devido a uma violação de bloqueio otimista
func IsOptimisticLockError(err error) bool {
	return errors.Is(err, ErrConcurrentModification)
}
