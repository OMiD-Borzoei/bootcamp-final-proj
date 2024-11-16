package models

type Voucher struct {
	Number  string `gorm:"primaryKey"`
	Version uint
	Items   []Voucheritem `gorm:"foreignKey:VoucherNumber"` // Specify the foreign key`
	//Items   []uint // store ids of voucherItems here
}

func (v *Voucher) ValidateVoucher() error {
	// Code Validation:
	if err := ValidateString(v.Number, "Number"); err != nil {
		return err
	}
	return nil
}
