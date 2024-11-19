package models

import "fmt"

type Voucher struct {
	ID      uint   `gorm:"primaryKey"`
	Number  string `gorm:"unique"` // Repetetive Number is handled in DB
	Version uint
	// a list of VIs that will be deleted when voucher gets deleted. since a vi can only belong to 1 voucher
	Items []Voucheritem `gorm:"foreignKey:VoucherID;constraint:onDelete:CASCADE"`
}

// non-DB Validation:
func (v *Voucher) Validate() error {
	//1- Number Validation:
	if err := ValidateString(v.Number, "Number"); err != nil {
		return err
	}
	//2- Number of vis Validation:
	if len(v.Items) > 500 || len(v.Items) < 2 {
		return fmt.Errorf("voucher has unacceptable number of voucheritems. expected between [2, 500], got %v", len(v.Items))
	}

	//3- Balance Validation:
	sum := 0
	for _, item := range v.Items {
		sum += int(item.Credit)
		sum -= int(item.Debit)
	}
	if sum != 0 {
		return fmt.Errorf("unbalance voucher")
	}

	return nil
}

// returns a voucher that is supposably Valid, or err.
func NewVoucher(number string) (*Voucher, error) {
	newV := &Voucher{
		Number: number,
	}
	if err := newV.Validate(); err != nil {
		return nil, err
	}

	return newV, nil
}
