package cbr

import (
	"test-service-b/internal/entity"

	"go.uber.org/zap"
)

type ICbrService interface {
	GetValCurs() (*[]entity.ValCurs)
}

type cbrService struct {
	repo   ICbrRepository
	logger *zap.Logger
}

func NewCbrService(logger *zap.Logger, repo ICbrRepository) ICbrService {
	cbrService := cbrService{
		repo: repo,
		logger: logger,
	}
	return cbrService
}

func (r cbrService) GetValCurs() (*[]entity.ValCurs) {
	return r.repo.GetValCurs()
}
