package saga

import (
	"context"
	"fmt"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelguudsPlatform/internal/merchant/service_errors"
	"go.uber.org/zap"
)

type SagaCoordinator struct {
	Logger *zap.Logger
}

// NewSagaCoordinator creates a new saga object
func NewSagaCoordinator(logger *zap.Logger) *SagaCoordinator {
	return &SagaCoordinator{Logger: logger}
}

// RunSaga creates a new saga instance and runs the required steps
func (s *SagaCoordinator) RunSaga(ctx context.Context, operationName string, steps ...*saga.Step) error {
	// define saga
	tx := saga.NewSaga(operationName)
	store := saga.New()

	for _, step := range steps {
		// first operation is to perform a distributed transaction and unlock the account if possible
		if err := tx.AddStep(step); err != nil {
			s.Logger.Error(service_errors.ErrFailedToConfigureSaga.Error())
			return service_errors.ErrFailedToConfigureSaga
		}
	}

	coordinator := saga.NewCoordinator(ctx, ctx, tx, store)
	if result := coordinator.Play(); result != nil && (len(result.CompensateErrors) > 0 || result.ExecutionError != nil) {
		// log the saga operation errors
		s.Logger.Error(service_errors.ErrSagaFailedToExecuteSuccessfully.Error(),
			zap.Errors("compensate error", result.CompensateErrors), zap.Error(result.ExecutionError))

		// construct error
		errMsg := fmt.Sprintf("compensate errors : %s , execution errors %s", zap.Errors("compensate error",
			result.CompensateErrors).String, zap.Error(result.ExecutionError).String)
		err := service_errors.NewError(errMsg)
		return err
	}

	return nil
}
