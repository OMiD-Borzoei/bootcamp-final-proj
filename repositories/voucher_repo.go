package repositories

import (
	"Project/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type VoucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) *VoucherRepository {
	return &VoucherRepository{db: db}
}

func (dr *VoucherRepository) Create(v *models.Voucher) (uint, error) {
	// 1- Validate the voucher itself first:
	if err := v.Validate(); err != nil {
		return 0, err
	}

	// 2- make an empty voucher with the given number and insert it in db:
	emptyV := models.Voucher{
		Number: v.Number,
	}
	if err := dr.db.Create(&emptyV).Error; err != nil {
		// check if unique constraint is violated:
		if strings.Contains(err.Error(), "23505") {
			return 0, fmt.Errorf("can't create Voucher cz a Voucher with number=%s already exists", emptyV.Number)
		}
		return 0, fmt.Errorf("could not create Voucher: %w", err)
	}

	// 3- insert all the items in db:
	viRepo := NewVoucherItemRepository(dr.db)
	for _, item := range v.Items {
		// Creating a voucheritem must first validate the itsem itself.
		if _, err := viRepo.Create(emptyV.ID, item.SLID, item.DLID, item.Debit, item.Credit); err != nil {
			// if any error happened delete the voucher made and return:
			dr.Delete(emptyV.ID)
			return 0, err
		}
	}

	return emptyV.ID, nil
}

func (dr *VoucherRepository) Read(id uint) (*models.Voucher, error) {

	var voucher models.Voucher
	// Load the Voucher along with its associated Voucheritems
	err := dr.db.Preload("Items").First(&voucher, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &voucher, nil
}

func (dr *VoucherRepository) ReadByNumber(number string) (*models.Voucher, error) {

	var voucher models.Voucher
	// Load the Voucher along with its associated Voucheritems
	err := dr.db.Preload("Items").First(&voucher, "number = ?", number).Error
	if err != nil {
		return nil, err
	}

	return &voucher, nil
}

func (dr *VoucherRepository) ReadAll() ([]models.Voucher, error) {
	var vList []models.Voucher
	result := dr.db.Preload("Items").Find(&vList)
	return vList, result.Error
}

func (dr *VoucherRepository) Update(id uint, v *models.Voucher) error {

	// 1- Ensure the given number is a valid string:
	if err := models.ValidateString(v.Number, "Number"); err != nil {
		return err
	}

	// 2- fetch the voucher that needs to be updated
	fetchedV := models.Voucher{}
	if err := dr.db.Model(&models.Voucher{}).Preload("Items").First(&fetchedV, "id = ?", id).Error; err != nil {
		return err
	}

	// 3- Handle Row Versioning:
	if fetchedV.Version != v.Version {
		return fmt.Errorf("version mismatch: expected %d but found %d", fetchedV.Version, v.Version)
	}

	// 4- Get the list of Inserted, Deleted and Updated vis:
	vis_partitioned := make(map[string][]*models.Voucheritem)
	// 5- See if this update will result in an unbalanced voucher:
	sum_debit, sum_credit := 0, 0
	viRepo := NewVoucherItemRepository(dr.db)

	// if the given vi is not already in db, this vi was meant to be INSERTED
	// else if the given vi has just an id field and the rest are zero-valued, this vi was meant to be DELETED
	// else the given vi was meant to be UPDATED
	for _, vi := range v.Items {

		vi.VoucherID = fetchedV.ID
		var fetchedVI models.Voucheritem

		if err := dr.db.Model(&models.Voucheritem{}).First(&fetchedVI, "id = ?", vi.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {

				if err1 := viRepo.ValidateVoucherItem(&vi); err1 != nil {
					return err1
				}

				vis_partitioned["toInsert"] = append(vis_partitioned["toInsert"], &vi)

				sum_credit += int(vi.Credit)
				sum_debit += int(vi.Debit)
			} else {
				return err
			}

		} else if vi.Credit == 0 && vi.Debit == 0 { //&& vi.SLID == 0 && vi.DLID == nil && vi.VoucherID == 0 {

			vis_partitioned["toDelete"] = append(vis_partitioned["toDelete"], &vi)

			sum_credit -= int(fetchedVI.Credit)
			sum_debit -= int(fetchedVI.Debit)

		} else {

			if err := viRepo.ValidateVoucherItem(&vi); err != nil {
				return err
			}

			vis_partitioned["toUpdate"] = append(vis_partitioned["toUpdate"], &vi)

			sum_credit += int(vi.Credit) - int(fetchedVI.Credit)
			sum_debit += int(vi.Debit) - int(fetchedVI.Debit)
		}

	}

	if sum_debit != sum_credit {
		return fmt.Errorf("requested update will result in an unbalanced voucher:\tsum_debit: %v\tsum_credit: %v", sum_debit, sum_credit)
	}

	newlen := len(fetchedV.Items) + len(vis_partitioned["toInsert"]) - len(vis_partitioned["toDelete"])
	if newlen < 2 || newlen > 500 {
		return fmt.Errorf("this update results in unacceptable number of voucheritems. expected between [2, 500], got %v", newlen)
	}

	// 6- Apply all changes, if any of them failed, rollback
	return dr.db.Transaction(func(tx *gorm.DB) error {
		newV := models.Voucher{
			ID:      id,
			Number:  v.Number,
			Version: v.Version + 1,
		}

		if err := dr.db.Model(&models.Voucher{}).Where("id = ?", id).Save(&newV).Error; err != nil {
			// check if unique constraint is violated:
			if strings.Contains(err.Error(), "23505") {
				return fmt.Errorf("can't update Voucher cz a Voucher with number=%s already exists", newV.Number)
			}
			return err
		}

		for _, vi := range vis_partitioned["toInsert"] {

			if _, err := viRepo.Create(vi.VoucherID, vi.SLID, vi.DLID, vi.Debit, vi.Credit); err != nil {
				return err
			}
		}

		for _, vi := range vis_partitioned["toUpdate"] {
			if err := viRepo.Update(vi.ID, vi); err != nil {
				return err
			}
		}

		for _, vi := range vis_partitioned["toDelete"] {
			if err := viRepo.Delete(vi.ID); err != nil {
				return nil
			}
		}

		// return nil to commit transaction if everything is fine:
		return nil
	})

}

func (dr *VoucherRepository) Delete(id uint) error {
	return dr.db.Model(&models.Voucher{}).Delete("id = ?", id).Error
}

func (dr *VoucherRepository) DeleteByNumber(number string) error {
	return dr.db.Model(&models.Voucher{}).Delete("number = ?", number).Error
}
