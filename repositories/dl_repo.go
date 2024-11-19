package repositories

import (
	"Project/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type DLRepository struct {
	db *gorm.DB
}

func NewDLRepository(db *gorm.DB) *DLRepository {
	return &DLRepository{db: db}
}

func (dr *DLRepository) Create(code string, title string) (uint, error) {

	dl, err := models.NewDL(code, title)
	if err != nil {
		return 0, err
	}

	if err := dr.db.Create(dl).Error; err != nil {
		// check if unique constraint is violated:
		if strings.Contains(err.Error(), "23505") {
			if strings.Contains(err.Error(), "code") {
				return 0, fmt.Errorf("can't create DL cz a DL with code=%s already exists", dl.Code)
			}
			if strings.Contains(err.Error(), "title") {
				return 0, fmt.Errorf("can't create DL cz a DL with title=%s already exists", dl.Title)
			}
		}
		return 0, fmt.Errorf("could not create DL: %w", err)
	}
	return dl.ID, nil
}

func (dr *DLRepository) Read(id uint) (*models.DL, error) {
	dl := &models.DL{}
	result := dr.db.Model(&models.DL{}).First(dl, "id = ?", id)

	return dl, result.Error
}

func (dr *DLRepository) ReadAll() ([]models.DL, error) {
	var dlList []models.DL
	result := dr.db.Find(&dlList)
	return dlList, result.Error
}

func (dr *DLRepository) ReadByCode(code string) (*models.DL, error) {
	// No need to use DB if code ain't Valid:
	if err := models.ValidateString(code, "Code"); err != nil {
		return nil, err
	}

	dl := &models.DL{}
	result := dr.db.Model(&models.DL{}).First(dl, "code = ?", code)

	return dl, result.Error
}

func (dr *DLRepository) ReadByTitle(title string) (*models.DL, error) {
	// No need to use DB if title ain't Valid:
	if err := models.ValidateString(title, "Title"); err != nil {
		return nil, err
	}

	dl := &models.DL{}
	result := dr.db.Model(&models.DL{}).First(dl, "title = ?", title)

	return dl, result.Error
}

func (dr *DLRepository) Update(id uint, newDL *models.DL) error {
	// 1- Check if newdl is valid:
	if err := newDL.Validate(); err != nil {
		return err
	}

	// 2- Check if There is a DL with such id in db:
	dl := &models.DL{}
	if result := dr.db.Model(&models.DL{}).First(dl, "id = ?", id); result.Error != nil {
		return result.Error
	}

	// 3- Handle Row Version:
	if dl.Version != newDL.Version {
		return fmt.Errorf("version mismatch: expected %d but found %d", dl.Version, newDL.Version)
	}

	// 4- Update:
	newDL.ID, newDL.Version = dl.ID, dl.Version+1
	err := dr.db.Model(&models.SL{}).Where("id = ?", id).Save(newDL).Error

	if err != nil {
		// check if unique constraint is violated:
		if strings.Contains(err.Error(), "23505") {
			if strings.Contains(err.Error(), "code") {
				return fmt.Errorf("can't update DL cz a DL with code=%s already exists", newDL.Code)
			}
			if strings.Contains(err.Error(), "title") {
				return fmt.Errorf("can't update DL cz a DL with title=%s already exists", newDL.Title)
			}
		}
		return fmt.Errorf("could not update DL: %w", err)
	}

	return nil
}

func (dr *DLRepository) Delete(id uint) error {
	err := dr.db.Model(&models.DL{}).Delete("id = ?", id).Error
	if err != nil && strings.Contains(err.Error(), "23503") {
		return fmt.Errorf("can't delete DL with id = %d cz it is referenced", id)
	}
	return err
}

func (dr *DLRepository) DeleteByCode(code string) error {
	// No need to use DB if code ain't Valid:
	if err := models.ValidateString(code, "Code"); err != nil {
		return err
	}

	err := dr.db.Model(&models.DL{}).Delete("code = ?", code).Error
	if err != nil && strings.Contains(err.Error(), "23503") {
		return fmt.Errorf("can't delete DL with code = %s cz it is referenced", code)
	}
	return err
}

func (dr *DLRepository) DeleteByTitle(title string) error {
	// No need to use DB if code ain't Valid:
	if err := models.ValidateString(title, "Title"); err != nil {
		return err
	}

	err := dr.db.Model(&models.DL{}).Delete("title = ?", title).Error
	if err != nil && strings.Contains(err.Error(), "23503") {
		return fmt.Errorf("can't delete DL with title = %s cz it is referenced", title)
	}
	return err
}
