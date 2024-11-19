package services

import (
	"Project/models"
	"Project/repositories"
)

type DLService struct {
	repo *repositories.DLRepository
}

func NewDLService(repo *repositories.DLRepository) *DLService {
	return &DLService{repo: repo}
}

func (s *DLService) GetAllDLs() ([]models.DL, error) {
	return s.repo.ReadAll()
}

func (s *DLService) GetDLByID(id uint) (*models.DL, error) {
	return s.repo.Read(id)
}

func (s *DLService) CreateDL(dl models.DL) (uint, error) {
	return s.repo.Create(dl.Code, dl.Title)
}

func (s *DLService) UpdateDL(id uint, dl models.DL) error {
	return s.repo.Update(id, &dl)
}

func (s *DLService) DeleteDL(id uint) error {
	return s.repo.Delete(id)
}
