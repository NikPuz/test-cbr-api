package cbr

import (
	"sync"
	"test-service-b/internal/entity"

	"go.uber.org/zap"
)

type ICbrRepository interface {
	GetValCurs() (*[]entity.ValCurs)
}

type cbrRepository struct {
	mu 				*sync.Mutex
	logger			*zap.Logger
	courseValues 	*[]entity.ValCurs 
}

func NewCbrRepository(mutex *sync.Mutex, logger *zap.Logger, courseValues *[]entity.ValCurs) ICbrRepository {
	cbrRepository := cbrRepository{
		mu: mutex,
		logger: logger,
		courseValues: courseValues,
	}
	return &cbrRepository
}

func (r *cbrRepository) GetValCurs() (*[]entity.ValCurs) {
	r.mu.Lock()
	result := *r.courseValues
	*r.courseValues = []entity.ValCurs{}
	r.mu.Unlock()
	return &result 
}
