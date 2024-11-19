package services

import (
	"Project/models"
	"Project/repositories"
)

type SLService struct {
	repo *repositories.SLRepository
}

func NewSLService(repo *repositories.SLRepository) *SLService {
	return &SLService{repo: repo}
}

func (s *SLService) CreateSL(code string, title string, hasDL bool) (uint, error) {
	return s.repo.Create(code, title, hasDL)
}

func (s *SLService) GetAllSLs() ([]models.SL, error) {
	return s.repo.ReadAll()
}

func (s *SLService) GetSLByID(id uint) (*models.SL, error) {
	return s.repo.Read(id)
}

func (s *SLService) GetSLByCode(code string) (*models.SL, error) {
	return s.repo.ReadByCode(code)
}

func (s *SLService) GetSLByTitle(title string) (*models.SL, error) {
	return s.repo.ReadByTitle(title)
}

func (s *SLService) UpdateSL(id uint, sl *models.SL) error {
	return s.repo.Update(id, sl)
}

func (s *SLService) DeleteSL(id uint) error {
	return s.repo.Delete(id)
}
