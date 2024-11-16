package repositories

import (
	"Project/models"
	"fmt"

	"gorm.io/gorm"
)

type SLRepository struct {
	db *gorm.DB
}

func NewSLRepository(db *gorm.DB) *SLRepository {
	return &SLRepository{db: db}
}

func (dr *SLRepository) Create(code string, title string, hassl bool) error {

	sl, err := models.NewSL(code, title, hassl)
	if err != nil {
		return err
	}

	if err := dr.db.Create(sl).Error; err != nil {
		return fmt.Errorf("could not create sl: %w", err)
	}
	return nil
}

func (dr *SLRepository) Read(id int) (*models.SL, error) {
	sl := &models.SL{}
	result := dr.db.Model(&models.SL{}).First(sl, "id = ?", id)

	return sl, result.Error
}

func (dr *SLRepository) ReadByCode(code string) (*models.SL, error) {
	if err := models.ValidateString(code, "Code"); err != nil {
		return nil, err
	}

	sl := &models.SL{}
	result := dr.db.Model(&models.SL{}).First(sl, "code = ?", code)

	return sl, result.Error
}

func (dr *SLRepository) ReadByTitle(title string) (*models.SL, error) {
	if err := models.ValidateString(title, "Title"); err != nil {
		return nil, err
	}

	sl := &models.SL{}
	result := dr.db.Model(&models.SL{}).First(sl, "title = ?", title)

	return sl, result.Error
}

func (dr *SLRepository) Update(id uint, newSL *models.SL) error {
	// 1- Check if newdl is valid:
	if err := newSL.Validate(); err != nil {
		return err
	}

	// 2- Check if There is a DL with such id in db:
	sl := &models.SL{}
	if result := dr.db.Model(&models.SL{}).First(sl, "id = ?", id); result.Error != nil {
		return result.Error
	}

	// 3- Handle Row Version:
	if sl.Version != newSL.Version {
		return fmt.Errorf("version mismatch: expected %d but found %d", sl.Version, newSL.Version)
	}

	// 4- Update
	newSL.ID = sl.ID
	newSL.Version = sl.Version + 1
	return dr.db.Model(&models.SL{}).Where("id = ?", id).Save(newSL).Error
}

func (dr *SLRepository) Delete(id uint) error {
	return dr.db.Model(&models.SL{}).Delete("id = ?", id).Error
}
