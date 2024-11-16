package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Voucheritem struct {
	gorm.Model
	VoucherNumber string
	SL            SL  `gorm:"foreignKey:Code;constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT;"`
	DL            *DL `gorm:"foreignKey:Code;constraint:OnDelete:RESTRICT;"`
	Debit         uint32
	Credit        uint32
}

func (vi *Voucheritem) ValidateVoucherItem() error {
	if vi.Debit > 0 && vi.Credit > 0 {
		return fmt.Errorf("one of the credit or debit must be 0")
	}
	return nil
}
