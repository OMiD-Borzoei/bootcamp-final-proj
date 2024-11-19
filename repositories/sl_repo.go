package repositories

import (
	"Project/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type SLRepository struct {
	db *gorm.DB
}

func NewSLRepository(db *gorm.DB) *SLRepository {
	return &SLRepository{db: db}
}

func (dr *SLRepository) Create(code string, title string, hassl bool) (uint, error) {

	sl, err := models.NewSL(code, title, hassl)
	if err != nil {
		return 0, err
	}

	if err := dr.db.Create(sl).Error; err != nil {
		return 0, fmt.Errorf("could not create sl: %w", err)
	}
	return sl.ID, nil
}

func (dr *SLRepository) Read(id uint) (*models.SL, error) {
	sl := &models.SL{}
	result := dr.db.Model(&models.SL{}).First(sl, "id = ?", id)

	return sl, result.Error
}

func (dr *SLRepository) ReadAll() ([]models.SL, error) {
	var slList []models.SL
	result := dr.db.Find(&slList)
	return slList, result.Error
}

func (dr *SLRepository) ReadByCode(code string) (*models.SL, error) {
	// No need to use DB if code ain't Valid:
	if err := models.ValidateString(code, "Code"); err != nil {
		return nil, err
	}

	sl := &models.SL{}
	result := dr.db.Model(&models.SL{}).First(sl, "code = ?", code)

	return sl, result.Error
}

func (dr *SLRepository) ReadByTitle(title string) (*models.SL, error) {
	// No need to use DB if title ain't Valid:
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

	// 2- Check if There is a SL with such id in db:
	sl := &models.SL{}
	if result := dr.db.Model(&models.SL{}).First(sl, "id = ?", id); result.Error != nil {
		return result.Error
	}

	// 3- Handle Row Version:
	if sl.Version != newSL.Version {
		return fmt.Errorf("version mismatch: expected %d but found %d", sl.Version, newSL.Version)
	}

	// 4- Check if there are any references to this sl:
	var count int64
	err := dr.db.Model(&models.Voucheritem{}).Where("sl_id = ?", sl.ID).Count(&count).Error
	// (SQLSTATE 42P01) : table voucheritem does not even exist
	if err != nil && !strings.Contains(err.Error(), "42P01") {
		return err
	}

	// If there are any references, raise an error
	if count > 0 {
		return fmt.Errorf("cannot update SL: it is referenced by %d Voucheritem(s)", count)
	}

	// 5- Update:
	newSL.ID = sl.ID
	newSL.Version = sl.Version + 1
	return dr.db.Model(&models.SL{}).Where("id = ?", id).Save(newSL).Error
}

func (dr *SLRepository) Delete(id uint) error {
	err := dr.db.Model(&models.SL{}).Delete("id = ?", id).Error
	if err != nil && strings.Contains(err.Error(), "23503") {
		return fmt.Errorf("can't delete SL with id = %d cz it is referenced", id)
	}
	return err
}

func (dr *SLRepository) DeleteByCode(code string) error {
	// No need to use DB if code ain't Valid:
	if err := models.ValidateString(code, "Code"); err != nil {
		return err
	}
	err := dr.db.Model(&models.SL{}).Delete("code = ?", code).Error
	if err != nil && strings.Contains(err.Error(), "23503") {
		return fmt.Errorf("can't delete SL with code = %s cz it is referenced", code)
	}
	return err
}

func (dr *SLRepository) DeleteByTitle(title string) error {
	// No need to use DB if title ain't Valid:
	if err := models.ValidateString(title, "Title"); err != nil {
		return err
	}
	err := dr.db.Model(&models.SL{}).Delete("title = ?", title).Error
	if err != nil && strings.Contains(err.Error(), "23503") {
		return fmt.Errorf("can't delete SL with title = %s cz it is referenced", title)
	}
	return err
}
