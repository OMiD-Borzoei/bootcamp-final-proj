package services

import (
	"Project/models"
	"Project/repositories"
	"fmt"
)

type VoucherService struct {
	repo *repositories.VoucherRepository
}

func NewVoucherService(repo *repositories.VoucherRepository) *VoucherService {
	return &VoucherService{repo: repo}
}

func (s *VoucherService) GetAllVouchers() ([]models.Voucher, error) {
	return s.repo.ReadAll()
}

func (vs *VoucherService) CreateVoucher(voucher *models.Voucher) (uint, error) {
	// Call repository method to create a new voucher
	voucherID, err := vs.repo.Create(voucher)
	if err != nil {
		return 0, fmt.Errorf("service: failed to create voucher: %w", err)
	}
	return voucherID, nil
}

func (vs *VoucherService) GetVoucher(id uint) (*models.Voucher, error) {
	// Call repository method to read a voucher
	voucher, err := vs.repo.Read(id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get voucher: %w", err)
	}
	return voucher, nil
}

func (vs *VoucherService) UpdateVoucher(id uint, voucher *models.Voucher) error {
	// Call repository method to update a voucher
	if err := vs.repo.Update(id, voucher); err != nil {
		return fmt.Errorf("service: failed to update voucher: %w", err)
	}
	return nil
}

func (vs *VoucherService) DeleteVoucher(id uint) error {
	// Call repository method to delete a voucher
	if err := vs.repo.Delete(id); err != nil {
		return fmt.Errorf("service: failed to delete voucher: %w", err)
	}
	return nil
}
