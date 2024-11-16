package repositories

import (
	"Project/models"
	"fmt"

	"gorm.io/gorm"
)

type DLRepository struct {
	db *gorm.DB
}

func NewDLRepository(db *gorm.DB) *DLRepository {
	return &DLRepository{db: db}
}

func (dr *DLRepository) Create(code string, title string) error {

	dl, err := models.NewDL(code, title)
	if err != nil {
		return err
	}

	if err := dr.db.Create(dl).Error; err != nil {
		return fmt.Errorf("could not create DL: %w", err)
	}
	return nil
}

func (dr *DLRepository) ReadByCode(code string) (*models.DL, error) {
	if err := models.ValidateString(code, "Code"); err != nil {
		return nil, err
	}

	dl := &models.DL{}
	result := dr.db.Model(&models.DL{}).First(dl, "code = ?", code)

	return dl, result.Error
}

func (dr *DLRepository) ReadByTitle(title string) (*models.DL, error) {
	if err := models.ValidateString(title, "Title"); err != nil {
		return nil, err
	}

	dl := &models.DL{}
	result := dr.db.Model(&models.DL{}).First(dl, "title = ?", title)

	return dl, result.Error
}
