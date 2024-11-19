package handlers

import (
	"Project/models"
	"Project/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type VoucherHandler struct {
	service *services.VoucherService
}

func NewVoucherHandler(service *services.VoucherService) *VoucherHandler {
	return &VoucherHandler{service: service}
}

func (h *VoucherHandler) GetAllVouchers(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Fetch DLs from the service
	dls, err := h.service.GetAllVouchers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the response as JSON
	if err := json.NewEncoder(w).Encode(dls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateVoucher handles POST requests to create a voucher
func (vh *VoucherHandler) CreateVoucher(w http.ResponseWriter, r *http.Request) {
	var voucher models.Voucher
	if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Call the service to create a voucher
	voucherID, err := vh.service.CreateVoucher(&voucher)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating voucher: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the created voucher ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uint{"voucher_id": voucherID})
}

// GetVoucher handles GET requests to fetch a voucher by ID
func (vh *VoucherHandler) GetVoucher(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Call the service to get the voucher
	voucher, err := vh.service.GetVoucher(uint(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching voucher: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the voucher data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(voucher)
}

// UpdateVoucher handles PUT requests to update a voucher
func (vh *VoucherHandler) UpdateVoucher(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var voucher models.Voucher
	if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
		fmt.Println(r.Body)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Call the service to update the voucher
	if err := vh.service.UpdateVoucher(uint(id), &voucher); err != nil {
		http.Error(w, fmt.Sprintf("Error updating voucher: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Voucher updated successfully"})
}

// DeleteVoucher handles DELETE requests to remove a voucher
func (vh *VoucherHandler) DeleteVoucher(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Call the service to delete the voucher
	if err := vh.service.DeleteVoucher(uint(id)); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting voucher: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Voucher deleted successfully"})
}
