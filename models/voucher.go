package models

import "fmt"

type Voucher struct {
	ID      uint   `gorm:"primaryKey"`
	Number  string `gorm:"unique"`
	Version uint
	Items   []Voucheritem `gorm:"foreignKey:VoucherID;constraint:onDelete:CASCADE"` // Specify the foreign key`
	//Items   []uint // store ids of voucherItems here
}

func (v *Voucher) Validate() error {
	// Code Validation:
	if err := ValidateString(v.Number, "Number"); err != nil {
		return err
	}

	if len(v.Items) > 500 || len(v.Items) < 2 {
		return fmt.Errorf("voucher has unacceptable number of voucheritems. expected between [2, 500], got %v", len(v.Items))
	}

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

func NewVoucher(number string) (*Voucher, error) {
	newV := &Voucher{
		Number: number,
	}
	if err := newV.Validate(); err != nil {
		return nil, err
	}

	return newV, nil
}
