package models

import (
	"fmt"
)

type Voucheritem struct {
	ID        uint `gorm:"primaryKey"`
	VoucherID uint
	SL        SL `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT;"`
	SLID      uint
	DL        *DL `gorm:"constraint:OnDelete:RESTRICT;"`
	DLID      *uint
	Debit     uint32
	Credit    uint32
}

func (vi *Voucheritem) Validate() error {
	if vi.Debit > 0 && vi.Credit > 0 {
		return fmt.Errorf("one of the credit or debit must be 0")
	}

	return nil
}

func NewVoucherItem(VoucherID uint, sl uint, dl *uint, debit, credit uint32) (*Voucheritem, error) {
	newvoucheritem := Voucheritem{
		VoucherID: VoucherID,
		SLID:      sl,
		DLID:      dl,
		Debit:     debit,
		Credit:    credit,
	}
	return &newvoucheritem, nil
}
