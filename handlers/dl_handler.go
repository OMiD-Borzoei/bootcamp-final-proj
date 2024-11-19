package handlers

import (
	"Project/models"
	"Project/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type DLHandler struct {
	service *services.DLService
}

func NewDLHandler(service *services.DLService) *DLHandler {
	return &DLHandler{service: service}
}

func (h *DLHandler) GetAllDLs(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Fetch DLs from the service
	dls, err := h.service.GetAllDLs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the response as JSON
	if err := json.NewEncoder(w).Encode(dls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *DLHandler) GetDLByID(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Parse the ID from the URL
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	dl, err := h.service.GetDLByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Encode and send the response as JSON
	if err := json.NewEncoder(w).Encode(dl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *DLHandler) CreateDL(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into the DL model
	var dl models.DL
	if err := json.NewDecoder(r.Body).Decode(&dl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the service to create the DL
	id, err := h.service.CreateDL(dl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the ID of the created DL
	response := map[string]interface{}{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *DLHandler) UpdateDL(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Parse the ID from the URL
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Decode the request body into the DL model
	var dl models.DL
	if err := json.NewDecoder(r.Body).Decode(&dl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the service to update the DL
	if err := h.service.UpdateDL(uint(id), dl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the updated DL ID
	response := map[string]interface{}{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *DLHandler) DeleteDL(w http.ResponseWriter, r *http.Request) {
	// Parse the ID from the URL
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Call the service to delete the DL
	if err := h.service.DeleteDL(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Set the response status to No Content (204)
	w.WriteHeader(http.StatusNoContent)
}
