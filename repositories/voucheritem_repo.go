package repositories

import (
	"Project/models"
	"fmt"

	"gorm.io/gorm"
)

type VoucherItemRepository struct {
	db *gorm.DB
}

func NewVoucherItemRepository(db *gorm.DB) *VoucherItemRepository {
	return &VoucherItemRepository{db: db}
}

func (dr *VoucherItemRepository) ValidateVoucherItem(vi *models.Voucheritem) error {
	if err := vi.Validate(); err != nil {
		return err
	}

	// Retrieve and check if Voucher exits:
	var foundvoucher models.Voucher
	if err := dr.db.Model(&models.Voucher{}).Where("id = ?", vi.VoucherID).First(&foundvoucher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("voucher with Number %v does not exist", vi.VoucherID)
		}
		return fmt.Errorf("error retrieving Voucher: %w", err)
	}

	// Retrieve and check if SL exists
	var foundsl models.SL
	if err := dr.db.Where("id = ?", vi.SLID).First(&foundsl).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("SL with ID %d does not exist", vi.SLID)
		}
		return fmt.Errorf("error retrieving SL: %w", err)
	}

	// Check if DL exists if SL.HasDL is true
	if foundsl.HasDL {
		if vi.DLID == nil {
			return fmt.Errorf("dl_id cannot be nil cz given sl must have a dl")
		}

		var dlExists bool
		if err := dr.db.Model(&models.DL{}).Select("count(*) > 0").Where("id = ?", *vi.DLID).Find(&dlExists).Error; err != nil {
			return fmt.Errorf("error checking DL existence: %w", err)
		}
		if !dlExists {
			return fmt.Errorf("DL with ID %d does not exist", *vi.DLID)
		}
		return nil
	}

	if vi.DLID != nil {
		return fmt.Errorf("given sl must not have a dl but a dl_id is provided")
	}

	return nil
}

func (dr *VoucherItemRepository) Create(voucherID uint, sl uint, dl *uint, debit, credit uint32) error {
	voucheritem, err := models.NewVoucherItem(voucherID, sl, dl, debit, credit)
	if err != nil {
		return err
	}

	if err := dr.ValidateVoucherItem(voucheritem); err != nil {
		return err
	}

	if err := dr.db.Create(voucheritem).Error; err != nil {
		return fmt.Errorf("could not create DL: %w", err)
	}
	return nil
}

func (dr *VoucherItemRepository) Read(id uint) (*models.Voucheritem, error) {
	vi := &models.Voucheritem{}
	return vi, dr.db.Model(&models.Voucheritem{}).First(vi, id).Error
}

// must only be called by voucherRepo.Update()
func (dr *VoucherItemRepository) Update(id uint, newVI *models.Voucheritem) error {

	// 1- Check if new vi is valid:
	if err := dr.ValidateVoucherItem(newVI); err != nil {
		return err
	}

	// 2- Check if There is a VI with such id in db:
	vi := &models.Voucheritem{}
	if result := dr.db.Model(&models.Voucheritem{}).First(vi, "id = ?", id); result.Error != nil {
		return result.Error
	}

	// 3- Update:
	newVI.ID = vi.ID
	return dr.db.Model(&models.Voucheritem{}).Where("id = ?", id).Save(newVI).Error
}

func (dr *VoucherItemRepository) Delete(id uint) error {
	return dr.db.Model(&models.Voucheritem{}).Delete("id = ?", id).Error
}
